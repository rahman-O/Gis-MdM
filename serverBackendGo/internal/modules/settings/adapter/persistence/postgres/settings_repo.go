package postgres

import (
	"context"
	"database/sql"
	"strings"

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
		       s.newdevicegroupid,
		       COALESCE(s.phonenumberformat,'+9 (999) 999-99-99'),
		       COALESCE(s.custompropertyname1,''), COALESCE(s.custompropertyname2,''),
		       COALESCE(s.custompropertyname3,''),
		       COALESCE(s.custommultiline1,false), COALESCE(s.custommultiline2,false),
		       COALESCE(s.custommultiline3,false),
		       COALESCE(s.customsend1,false), COALESCE(s.customsend2,false),
		       COALESCE(s.customsend3,false),
		       COALESCE(s.desktopheadertemplate,''), COALESCE(s.senddescription,false),
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
	var newCfg, newGroup sql.NullInt64
	var bgURL sql.NullString
	if err := row.Scan(
		&s.ID, &s.CustomerID, &s.CustomerName, &s.TwoFactor, &idle,
		&s.Language, &s.UseDefaultLanguage, &s.CreateNewDevices, &newCfg,
		&newGroup, &s.PhoneNumberFormat, &s.CustomPropertyName1, &s.CustomPropertyName2,
		&s.CustomPropertyName3, &s.CustomMultiline1, &s.CustomMultiline2, &s.CustomMultiline3,
		&s.CustomSend1, &s.CustomSend2, &s.CustomSend3, &s.DesktopHeaderTemplate, &s.SendDescription,
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
	if newGroup.Valid {
		v := int(newGroup.Int64)
		s.NewDeviceGroupID = &v
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
		PasswordLength: 6, PhoneNumberFormat: "+9 (999) 999-99-99",
		BackgroundColor: "#678ca6", TextColor: "#ffffff",
		IconSize: "SMALL", DesktopHeader: "NO_HEADER", SingleCustomer: true,
	}
}

func (r *SettingsRepository) SaveMisc(ctx context.Context, s *domain.Settings) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE settings SET
			createnewdevices=$1, newdeviceconfigurationid=$2,
			newdevicegroupid=$3, phonenumberformat=$4,
			custompropertyname1=$5, custompropertyname2=$6, custompropertyname3=$7,
			custommultiline1=$8, custommultiline2=$9, custommultiline3=$10,
			customsend1=$11, customsend2=$12, customsend3=$13,
			desktopheadertemplate=$14, senddescription=$15,
			passwordlength=$16, passwordstrength=$17, twofactor=$18, idlelogout=$19
		WHERE customerid=$20`,
		s.CreateNewDevices, nullInt(s.NewDeviceConfigurationID), nullInt(s.NewDeviceGroupID),
		nullStrVal(s.PhoneNumberFormat), nullStrVal(s.CustomPropertyName1),
		nullStrVal(s.CustomPropertyName2), nullStrVal(s.CustomPropertyName3),
		s.CustomMultiline1, s.CustomMultiline2, s.CustomMultiline3,
		s.CustomSend1, s.CustomSend2, s.CustomSend3,
		nullStrVal(s.DesktopHeaderTemplate), s.SendDescription,
		s.PasswordLength, s.PasswordStrength, s.TwoFactor, nullIntPtr(s.IdleLogout), s.CustomerID)
	return err
}

func nullStrVal(s string) interface{} {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return s
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

const userRoleSelect = `
	SELECT roleid, customerid,
		columndisplayeddevicestatus, columndisplayeddevicedate, columndisplayeddevicenumber,
		columndisplayeddevicemodel, columndisplayeddevicepermissionsstatus,
		columndisplayeddeviceappinstallstatus, columndisplayeddeviceconfiguration,
		columndisplayeddeviceimei, columndisplayeddevicephone, columndisplayeddevicedesc,
		columndisplayeddevicegroup, columndisplayedlauncherversion,
		columndisplayeddevicefilesstatus, columndisplayedbatterylevel,
		columndisplayeddefaultlauncher, columndisplayedcustom1, columndisplayedcustom2,
		columndisplayedcustom3, columndisplayedmdmmode, columndisplayedkioskmode,
		columndisplayedandroidversion, columndisplayedenrollmentdate,
		columndisplayedserial, columndisplayedpublicip
	FROM userrolesettings
	WHERE customerid = $1 AND roleid = $2`

func (r *SettingsRepository) GetUserRoleSettings(ctx context.Context, customerID, roleID int) (*domain.UserRoleSettings, error) {
	var s domain.UserRoleSettings
	err := r.db.QueryRowContext(ctx, userRoleSelect, customerID, roleID).Scan(
		&s.RoleID, &s.CustomerID,
		&s.ColumnDisplayedDeviceStatus, &s.ColumnDisplayedDeviceDate, &s.ColumnDisplayedDeviceNumber,
		&s.ColumnDisplayedDeviceModel, &s.ColumnDisplayedDevicePermissionsStatus,
		&s.ColumnDisplayedDeviceAppInstallStatus, &s.ColumnDisplayedDeviceConfiguration,
		&s.ColumnDisplayedDeviceImei, &s.ColumnDisplayedDevicePhone, &s.ColumnDisplayedDeviceDesc,
		&s.ColumnDisplayedDeviceGroup, &s.ColumnDisplayedLauncherVersion,
		&s.ColumnDisplayedDeviceFilesStatus, &s.ColumnDisplayedBatteryLevel,
		&s.ColumnDisplayedDefaultLauncher, &s.ColumnDisplayedCustom1, &s.ColumnDisplayedCustom2,
		&s.ColumnDisplayedCustom3, &s.ColumnDisplayedMdmMode, &s.ColumnDisplayedKioskMode,
		&s.ColumnDisplayedAndroidVersion, &s.ColumnDisplayedEnrollmentDate,
		&s.ColumnDisplayedSerial, &s.ColumnDisplayedPublicIp,
	)
	if err == sql.ErrNoRows {
		def := domain.DefaultUserRoleSettings(roleID, customerID)
		return &def, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SettingsRepository) SaveUserRoleSettings(ctx context.Context, customerID int, s domain.UserRoleSettings) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO userrolesettings (
			roleid, customerid,
			columndisplayeddevicestatus, columndisplayeddevicedate, columndisplayeddevicenumber,
			columndisplayeddevicemodel, columndisplayeddevicepermissionsstatus,
			columndisplayeddeviceappinstallstatus, columndisplayeddeviceconfiguration,
			columndisplayeddeviceimei, columndisplayeddevicephone, columndisplayeddevicedesc,
			columndisplayeddevicegroup, columndisplayedlauncherversion,
			columndisplayeddevicefilesstatus, columndisplayedbatterylevel,
			columndisplayeddefaultlauncher, columndisplayedcustom1, columndisplayedcustom2,
			columndisplayedcustom3, columndisplayedmdmmode, columndisplayedkioskmode,
			columndisplayedandroidversion, columndisplayedenrollmentdate,
			columndisplayedserial, columndisplayedpublicip
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26)
		ON CONFLICT (roleid, customerid) DO UPDATE SET
			columndisplayeddevicestatus=$3, columndisplayeddevicedate=$4, columndisplayeddevicenumber=$5,
			columndisplayeddevicemodel=$6, columndisplayeddevicepermissionsstatus=$7,
			columndisplayeddeviceappinstallstatus=$8, columndisplayeddeviceconfiguration=$9,
			columndisplayeddeviceimei=$10, columndisplayeddevicephone=$11, columndisplayeddevicedesc=$12,
			columndisplayeddevicegroup=$13, columndisplayedlauncherversion=$14,
			columndisplayeddevicefilesstatus=$15, columndisplayedbatterylevel=$16,
			columndisplayeddefaultlauncher=$17, columndisplayedcustom1=$18, columndisplayedcustom2=$19,
			columndisplayedcustom3=$20, columndisplayedmdmmode=$21, columndisplayedkioskmode=$22,
			columndisplayedandroidversion=$23, columndisplayedenrollmentdate=$24,
			columndisplayedserial=$25, columndisplayedpublicip=$26`,
		s.RoleID, customerID,
		s.ColumnDisplayedDeviceStatus, s.ColumnDisplayedDeviceDate, s.ColumnDisplayedDeviceNumber,
		s.ColumnDisplayedDeviceModel, s.ColumnDisplayedDevicePermissionsStatus,
		s.ColumnDisplayedDeviceAppInstallStatus, s.ColumnDisplayedDeviceConfiguration,
		s.ColumnDisplayedDeviceImei, s.ColumnDisplayedDevicePhone, s.ColumnDisplayedDeviceDesc,
		s.ColumnDisplayedDeviceGroup, s.ColumnDisplayedLauncherVersion,
		s.ColumnDisplayedDeviceFilesStatus, s.ColumnDisplayedBatteryLevel,
		s.ColumnDisplayedDefaultLauncher, s.ColumnDisplayedCustom1, s.ColumnDisplayedCustom2,
		s.ColumnDisplayedCustom3, s.ColumnDisplayedMdmMode, s.ColumnDisplayedKioskMode,
		s.ColumnDisplayedAndroidVersion, s.ColumnDisplayedEnrollmentDate,
		s.ColumnDisplayedSerial, s.ColumnDisplayedPublicIp,
	)
	return err
}
