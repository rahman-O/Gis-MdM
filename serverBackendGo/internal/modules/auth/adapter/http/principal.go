package http

import (
	"github.com/gis-mdm/server-backend-go/internal/modules/auth/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

func principalFromView(view *domain.UserView) *platformauth.Principal {
	if view == nil {
		return nil
	}
	perms := make([]string, 0)
	roleID := 0
	if view.UserRole != nil {
		roleID = view.UserRole.ID
		for _, p := range view.UserRole.Permissions {
			if p.Name != "" {
				perms = append(perms, p.Name)
			}
		}
	}
	return &platformauth.Principal{
		ID:            view.ID,
		Login:         view.Login,
		AuthToken:     view.AuthToken,
		CustomerID:    view.CustomerID,
		RoleID:        roleID,
		SuperAdmin:    view.SuperAdmin,
		Permissions:   perms,
		PasswordReset: view.PasswordReset,
		AuthLoaded:    true,
	}
}
