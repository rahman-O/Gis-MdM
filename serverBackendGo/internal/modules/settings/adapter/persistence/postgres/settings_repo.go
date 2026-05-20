package postgres

import (
	"context"
	"database/sql"

	"github.com/gis-mdm/server-backend-go/internal/modules/settings/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/settings/port"
)

type SettingsRepository struct {
	db *sql.DB
}

func NewSettingsRepository(db *sql.DB) *SettingsRepository {
	return &SettingsRepository{db: db}
}

var _ port.SettingsRepository = (*SettingsRepository)(nil)

func (r *SettingsRepository) IsSingleCustomer(ctx context.Context) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS (SELECT 1 FROM customers WHERE id > 1 LIMIT 1)`).Scan(&exists)
	return !exists, err
}

func (r *SettingsRepository) GetByCustomerID(ctx context.Context, customerID int) (*domain.Settings, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT s.id, s.customerid, COALESCE(c.name,''), s.twofactor, s.idlelogout,
		       COALESCE(s.language,'en'), COALESCE(s.usedefaultlanguage,true),
		       COALESCE(s.createnewdevices,false), s.newdeviceconfigurationid,
		       COALESCE(s.passwordlength,6), COALESCE(s.passwordstrength,0),
		       COALESCE(s.backgroundcolor,'#678ca6'), COALESCE(s.textcolor,'#ffffff'),
		       COALESCE(s.backgroundimageurl,''), COALESCE(s.iconsize,'SMALL'),
		       COALESCE(s.desktopheader,'NO_HEADER')
		FROM settings s
		INNER JOIN customers c ON c.id = s.customerid
		WHERE s.customerid = $1
		LIMIT 1`, customerID)

	var s domain.Settings
	var idle sql.NullInt64
	var newCfg sql.NullInt64
	var bgURL sql.NullString
	if err := row.Scan(
		&s.ID, &s.CustomerID, &s.CustomerName, &s.TwoFactor, &idle,
		&s.Language, &s.UseDefaultLanguage, &s.CreateNewDevices, &newCfg,
		&s.PasswordLength, &s.PasswordStrength,
		&s.BackgroundColor, &s.TextColor, &bgURL, &s.IconSize, &s.DesktopHeader,
	); err == sql.ErrNoRows {
		return r.defaultSettings(customerID), nil
	} else if err != nil {
		// schema without 000003 columns — minimal row
		return r.getMinimal(ctx, customerID)
	}
	if idle.Valid {
		v := int(idle.Int64)
		s.IdleLogout = &v
	}
	if newCfg.Valid {
		v := int(newCfg.Int64)
		s.NewDeviceConfigurationID = &v
	}
	if bgURL.Valid {
		s.BackgroundImageURL = bgURL.String
	}
	single, _ := r.IsSingleCustomer(ctx)
	s.SingleCustomer = single
	return &s, nil
}

func (r *SettingsRepository) getMinimal(ctx context.Context, customerID int) (*domain.Settings, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT s.id, s.customerid, COALESCE(c.name,''), s.twofactor, s.idlelogout
		FROM settings s
		INNER JOIN customers c ON c.id = s.customerid
		WHERE s.customerid = $1 LIMIT 1`, customerID)
	var s domain.Settings
	var idle sql.NullInt64
	err := row.Scan(&s.ID, &s.CustomerID, &s.CustomerName, &s.TwoFactor, &idle)
	if err == sql.ErrNoRows {
		return r.defaultSettings(customerID), nil
	}
	if err != nil {
		return nil, err
	}
	if idle.Valid {
		v := int(idle.Int64)
		s.IdleLogout = &v
	}
	s.Language = "en"
	s.UseDefaultLanguage = true
	s.PasswordLength = 6
	s.BackgroundColor = "#678ca6"
	s.TextColor = "#ffffff"
	s.IconSize = "SMALL"
	s.DesktopHeader = "NO_HEADER"
	single, _ := r.IsSingleCustomer(ctx)
	s.SingleCustomer = single
	return &s, nil
}

func (r *SettingsRepository) defaultSettings(customerID int) *domain.Settings {
	return &domain.Settings{
		ID: 0, CustomerID: customerID, Language: "en", UseDefaultLanguage: true,
		PasswordLength: 6, BackgroundColor: "#678ca6", TextColor: "#ffffff",
		IconSize: "SMALL", DesktopHeader: "NO_HEADER", SingleCustomer: true,
	}
}

func (r *SettingsRepository) SaveMisc(ctx context.Context, s *domain.Settings) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE settings SET
			createnewdevices=$1, newdeviceconfigurationid=$2,
			passwordlength=$3, passwordstrength=$4, twofactor=$5, idlelogout=$6
		WHERE customerid=$7`,
		s.CreateNewDevices, nullInt(s.NewDeviceConfigurationID),
		s.PasswordLength, s.PasswordStrength, s.TwoFactor, nullIntPtr(s.IdleLogout), s.CustomerID)
	return err
}

func (r *SettingsRepository) SaveLanguage(ctx context.Context, s *domain.Settings) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE settings SET language=$1, usedefaultlanguage=$2 WHERE customerid=$3`,
		s.Language, s.UseDefaultLanguage, s.CustomerID)
	return err
}

func (r *SettingsRepository) SaveDesign(ctx context.Context, s *domain.Settings) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE settings SET backgroundcolor=$1, textcolor=$2, backgroundimageurl=$3,
			iconsize=$4, desktopheader=$5 WHERE customerid=$6`,
		s.BackgroundColor, s.TextColor, s.BackgroundImageURL, s.IconSize, s.DesktopHeader, s.CustomerID)
	return err
}

func nullInt(p *int) interface{} {
	if p == nil {
		return nil
	}
	return *p
}

func nullIntPtr(p *int) interface{} {
	if p == nil {
		return nil
	}
	return *p
}
