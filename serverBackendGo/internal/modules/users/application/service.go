package application

import (
	"context"
	"errors"
	"strings"

	authdomain "github.com/gis-mdm/server-backend-go/internal/modules/auth/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/users/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/users/port"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
	"github.com/gis-mdm/server-backend-go/internal/shared/crypto"
)

// Service implements users use cases.
type Service struct {
	repo port.Repository
}

func NewService(repo port.Repository) *Service {
	return &Service{repo: repo}
}

var (
	ErrPermissionDenied = errors.New("error.permission.denied")
	ErrWrongPassword    = errors.New("error.password.wrong")
	ErrEmptyPassword    = errors.New("error.password.empty")
	ErrDuplicateEmail   = errors.New("error.duplicate.email")
	ErrDuplicateLogin   = errors.New("error.duplicate.login")
	ErrUserNotFound     = errors.New("error.user.not.found")
)

// GetCurrentUser returns the full user view for the session.
func (s *Service) GetCurrentUser(ctx context.Context, userID int64) (*authdomain.UserView, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return enrichAndView(ctx, s.repo, user)
}

func enrichAndView(ctx context.Context, repo port.Repository, user *authdomain.User) (*authdomain.UserView, error) {
	if user == nil {
		return nil, nil
	}
	settings, _ := repo.GetCustomerSettings(ctx, user.CustomerID)
	if settings != nil {
		if settings.TwoFactor {
			user.TwoFactor = true
		}
		user.IdleLogout = settings.IdleLogout
	}
	single, _ := repo.IsSingleCustomer(ctx)
	user.SingleCustomer = single
	user.Password = ""
	user.AuthToken = ""
	return authdomain.NewUserView(user), nil
}

// UpdateProfile updates name/email for the authenticated user.
func (s *Service) UpdateProfile(ctx context.Context, principal *platformauth.Principal, p domain.ProfilePayload) (*authdomain.UserView, error) {
	if principal == nil {
		return nil, ErrPermissionDenied
	}
	dbUser, err := s.repo.FindByID(ctx, principal.ID)
	if err != nil {
		return nil, err
	}
	if dbUser == nil {
		return nil, ErrUserNotFound
	}
	email := strings.TrimSpace(p.Email)
	if email == "" {
		dbUser.Email = ""
	} else if !strings.EqualFold(email, dbUser.Email) {
		other, err := s.repo.FindByEmail(ctx, email)
		if err != nil {
			return nil, err
		}
		if other != nil && other.ID != dbUser.ID {
			return nil, ErrDuplicateEmail
		}
		dbUser.Email = email
	}
	dbUser.Name = p.Name
	if err := s.repo.UpdateMainDetails(ctx, dbUser); err != nil {
		return nil, err
	}
	return enrichAndView(ctx, s.repo, dbUser)
}

// ChangePassword updates password for the authenticated user.
func (s *Service) ChangePassword(ctx context.Context, principal *platformauth.Principal, p domain.UserPayload) error {
	if principal == nil || p.ID == nil || *p.ID != principal.ID {
		return ErrPermissionDenied
	}
	if strings.TrimSpace(p.NewPassword) == "" {
		return ErrEmptyPassword
	}
	dbUser, err := s.repo.FindByLogin(ctx, p.Login)
	if err != nil {
		return err
	}
	if dbUser == nil {
		return ErrUserNotFound
	}
	if !crypto.PasswordMatch(p.OldPassword, dbUser.Password) {
		return ErrWrongPassword
	}
	hash := crypto.HashFromMd5(p.NewPassword)
	token := crypto.GenerateAuthToken()
	return s.repo.UpdatePassword(ctx, dbUser.ID, hash, token, true, false, nil)
}

// ListUsers returns tenant users for admin UI.
func (s *Service) ListUsers(ctx context.Context, principal *platformauth.Principal, filter string) ([]*authdomain.UserView, error) {
	if principal == nil || !principal.CanListUsers() {
		return nil, ErrPermissionDenied
	}
	users, err := s.repo.ListByCustomer(ctx, principal.CustomerID, filter)
	if err != nil {
		return nil, err
	}
	out := make([]*authdomain.UserView, 0, len(users))
	for _, u := range users {
		u.Password = ""
		u.AuthToken = ""
		v := authdomain.NewUserView(u)
		v.Editable = u.ID != principal.ID
		out = append(out, v)
	}
	return out, nil
}

// UpsertUser creates or updates a tenant user.
func (s *Service) UpsertUser(ctx context.Context, principal *platformauth.Principal, p domain.UserPayload) error {
	if principal == nil || !principal.CanManageUsers() {
		return ErrPermissionDenied
	}
	email := strings.TrimSpace(p.Email)
	if email != "" {
		other, err := s.repo.FindByEmail(ctx, email)
		if err != nil {
			return err
		}
		if other != nil && (p.ID == nil || other.ID != *p.ID) {
			return ErrDuplicateEmail
		}
	}
	byLogin, err := s.repo.FindByLogin(ctx, p.Login)
	if err != nil {
		return err
	}
	if byLogin != nil && (p.ID == nil || byLogin.ID != *p.ID) {
		return ErrDuplicateLogin
	}
	roleID := 3
	if p.UserRole != nil && p.UserRole.ID > 0 {
		roleID = p.UserRole.ID
	}
	allDev := true
	if p.AllDevicesAvailable != nil {
		allDev = *p.AllDevicesAvailable
	}
	allCfg := true
	if p.AllConfigAvailable != nil {
		allCfg = *p.AllConfigAvailable
	}
	if p.ID == nil {
		if strings.TrimSpace(p.NewPassword) == "" {
			return ErrEmptyPassword
		}
		u := &authdomain.User{
			Login:               p.Login,
			Name:                p.Name,
			Email:               email,
			Password:            crypto.HashFromMd5(p.NewPassword),
			CustomerID:          principal.CustomerID,
			AuthToken:           crypto.GenerateAuthToken(),
			AllDevicesAvailable: allDev,
			AllConfigAvailable:  allCfg,
			UserRole:            &authdomain.UserRole{ID: roleID},
			Groups:              p.Groups,
			Configurations:      p.Configurations,
		}
		reset, _ := s.repo.PasswordResetEnabled(ctx, principal.CustomerID)
		if reset {
			u.PasswordReset = true
			u.PasswordResetToken = crypto.GenerateAuthToken()
		}
		return s.repo.Insert(ctx, u)
	}
	dbUser, err := s.repo.FindByID(ctx, *p.ID)
	if err != nil {
		return err
	}
	if dbUser == nil || dbUser.CustomerID != principal.CustomerID {
		return ErrUserNotFound
	}
	dbUser.Login = p.Login
	dbUser.Name = p.Name
	dbUser.Email = email
	dbUser.AllDevicesAvailable = allDev
	dbUser.AllConfigAvailable = allCfg
	dbUser.Groups = p.Groups
	dbUser.Configurations = p.Configurations
	if p.UserRole != nil {
		dbUser.UserRole = &authdomain.UserRole{ID: roleID}
	}
	if err := s.repo.UpdateMainDetails(ctx, dbUser); err != nil {
		return err
	}
	if strings.TrimSpace(p.NewPassword) != "" {
		hash := crypto.HashFromMd5(p.NewPassword)
		token := crypto.GenerateAuthToken()
		reset, _ := s.repo.PasswordResetEnabled(ctx, principal.CustomerID)
		var resetTok *string
		passReset := false
		if reset {
			passReset = true
			t := crypto.GenerateAuthToken()
			resetTok = &t
		}
		return s.repo.UpdatePassword(ctx, dbUser.ID, hash, token, true, passReset, resetTok)
	}
	return nil
}

// DeleteUser removes another user in the tenant.
func (s *Service) DeleteUser(ctx context.Context, principal *platformauth.Principal, targetID int64) error {
	if principal == nil || !principal.CanManageUsers() {
		return ErrPermissionDenied
	}
	if principal.ID == targetID {
		return errors.New("cannot delete self")
	}
	target, err := s.repo.FindByID(ctx, targetID)
	if err != nil {
		return err
	}
	if target == nil || target.CustomerID != principal.CustomerID {
		return ErrUserNotFound
	}
	return s.repo.Delete(ctx, targetID)
}
