package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gis-mdm/server-backend-go/internal/modules/enrollment_routes/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/enrollment_routes/port"
)

type RouteRepository struct {
	db *sql.DB
}

func NewRouteRepository(db *sql.DB) *RouteRepository {
	return &RouteRepository{db: db}
}

var _ port.RouteRepository = (*RouteRepository)(nil)

var (
	ErrNotFound      = errors.New("enrollment route not found")
	ErrDuplicateName = errors.New("duplicate enrollment route name")
)

func (r *RouteRepository) ListViews(ctx context.Context, customerID int) ([]domain.EnrollmentRouteView, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT er.id, er.name, COALESCE(er.description, ''), COALESCE(er.qrcodekey, ''),
		       er.default_tree_node_id, COALESCE(n.name, ''), COALESCE(n.path, ''),
		       er.default_device_id_mode,
		       COALESCE(er.bootstrap_intent, 'stable'),
		       COALESCE(er.bootstrap_application_id, 0),
		       COALESCE(a.name, ''),
		       er.bootstrap_version_id,
		       er.mainappid,
		       COALESCE(av.version, ''),
		       COALESCE(a.pkg, ''),
		       er.container_placement_ack_at,
		       er.type,
		       er.wifi_ssid, er.wifi_password, er.wifi_security_type,
		       er.qr_parameters, er.admin_extras,
		       er.mobile_enrollment, er.encrypt_device
		FROM enrollment_routes er
		LEFT JOIN device_tree_nodes n ON n.id = er.default_tree_node_id
		LEFT JOIN applications a ON a.id = er.bootstrap_application_id
		LEFT JOIN applicationversions av ON av.id = er.mainappid
		WHERE er.customerid = $1
		ORDER BY lower(er.name)`, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.EnrollmentRouteView
	for rows.Next() {
		v, err := scanView(rows)
		if err != nil {
			return nil, err
		}
		if err := r.enrichView(ctx, customerID, &v); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

func (r *RouteRepository) GetViewByID(ctx context.Context, customerID, id int) (*domain.EnrollmentRouteView, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT er.id, er.name, COALESCE(er.description, ''), COALESCE(er.qrcodekey, ''),
		       er.default_tree_node_id, COALESCE(n.name, ''), COALESCE(n.path, ''),
		       er.default_device_id_mode,
		       COALESCE(er.bootstrap_intent, 'stable'),
		       COALESCE(er.bootstrap_application_id, 0),
		       COALESCE(a.name, ''),
		       er.bootstrap_version_id,
		       er.mainappid,
		       COALESCE(av.version, ''),
		       COALESCE(a.pkg, ''),
		       er.container_placement_ack_at,
		       er.type,
		       er.wifi_ssid, er.wifi_password, er.wifi_security_type,
		       er.qr_parameters, er.admin_extras,
		       er.mobile_enrollment, er.encrypt_device
		FROM enrollment_routes er
		LEFT JOIN device_tree_nodes n ON n.id = er.default_tree_node_id
		LEFT JOIN applications a ON a.id = er.bootstrap_application_id
		LEFT JOIN applicationversions av ON av.id = er.mainappid
		WHERE er.id = $1 AND er.customerid = $2`, id, customerID)
	v, err := scanView(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if err := r.enrichView(ctx, customerID, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

type scannable interface {
	Scan(dest ...any) error
}

func scanView(row scannable) (domain.EnrollmentRouteView, error) {
	var v domain.EnrollmentRouteView
	var treeID sql.NullInt64
	var bootAppID sql.NullInt64
	var bootVerID, mainApp sql.NullInt64
	var ackAt sql.NullTime
	var wifiSSID, wifiPassword, wifiSecurityType, qrParameters, adminExtras sql.NullString
	var mobileEnrollment, encryptDevice sql.NullBool
	if err := row.Scan(
		&v.ID, &v.Name, &v.Description, &v.QRCodeKey,
		&treeID, &v.TargetNodeName, &v.TargetNodePath,
		&v.DeviceIdentityMode,
		&v.BootstrapIntent, &bootAppID, &v.BootstrapApplicationName,
		&bootVerID, &mainApp, &v.ResolvedVersionLabel, &v.ResolvedPackage,
		&ackAt, &v.Type,
		&wifiSSID, &wifiPassword, &wifiSecurityType,
		&qrParameters, &adminExtras,
		&mobileEnrollment, &encryptDevice,
	); err != nil {
		return v, err
	}
	if wifiSSID.Valid {
		v.WifiSSID = wifiSSID.String
	}
	if wifiPassword.Valid {
		v.WifiPassword = wifiPassword.String
	}
	if wifiSecurityType.Valid {
		v.WifiSecurityType = wifiSecurityType.String
	}
	if qrParameters.Valid {
		v.QRParameters = qrParameters.String
	}
	if adminExtras.Valid {
		v.AdminExtras = adminExtras.String
	}
	if mobileEnrollment.Valid {
		v.MobileEnrollment = mobileEnrollment.Bool
	}
	if encryptDevice.Valid {
		v.EncryptDevice = encryptDevice.Bool
	}
	if treeID.Valid {
		v.TargetNodeID = int(treeID.Int64)
	}
	if bootAppID.Valid {
		v.BootstrapApplicationID = int(bootAppID.Int64)
	}
	if bootVerID.Valid && bootVerID.Int64 > 0 {
		b := int(bootVerID.Int64)
		v.BootstrapVersionID = &b
	}
	if mainApp.Valid && mainApp.Int64 > 0 {
		m := int(mainApp.Int64)
		v.ResolvedMainAppVersionID = &m
	}
	if ackAt.Valid {
		t := ackAt.Time
		v.ContainerPlacementAckAt = &t
		v.ContainerPlacementAcknowledged = true
	}
	if strings.TrimSpace(v.QRCodeKey) != "" {
		v.Status = "active"
	} else {
		v.Status = "draft"
	}
	return v, nil
}

func (r *RouteRepository) enrichView(ctx context.Context, customerID int, v *domain.EnrollmentRouteView) error {
	if v.TargetNodeID > 0 {
		kind, err := r.NodePlacementKind(ctx, customerID, v.TargetNodeID)
		if err != nil {
			return err
		}
		v.TargetPlacementKind = kind
	}
	return nil
}

func (r *RouteRepository) Create(ctx context.Context, customerID int, req domain.CreateRequest, qrcodeKey string, resolved domain.ResolvedBootstrap, containerAck bool) (int, error) {
	mode := deviceMode(req.DeviceIdentityMode, req.DefaultDeviceIDMode)
	typ := 0
	if req.Type != nil {
		typ = *req.Type
	}
	var ack any
	if containerAck {
		ack = time.Now()
	}
	var id int
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO enrollment_routes (
			customerid, name, description, qrcodekey, mainappid,
			profile_version_id, default_tree_node_id, default_device_id_mode, type,
			bootstrap_intent, bootstrap_application_id, bootstrap_version_id, container_placement_ack_at,
			wifi_ssid, wifi_password, wifi_security_type, qr_parameters, admin_extras,
			mobile_enrollment, encrypt_device
		) VALUES ($1,$2,$3,$4,$5,NULL,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19)
		RETURNING id`,
		customerID, strings.TrimSpace(req.Name), nullStr(req.Description), qrcodeKey, resolved.VersionID,
		req.TargetNodeID, mode, typ,
		req.BootstrapIntent, req.BootstrapApplicationID, nullIntPtr(req.BootstrapVersionID), ack,
		nullStr(req.WifiSSID), nullStr(req.WifiPassword), nullStr(req.WifiSecurityType),
		nullStr(req.QRParameters), nullStr(req.AdminExtras),
		nullBoolPtr(req.MobileEnrollment), nullBoolPtr(req.EncryptDevice),
	).Scan(&id)
	if err != nil && strings.Contains(err.Error(), "enrollment_routes_name_customer_uidx") {
		return 0, ErrDuplicateName
	}
	if err != nil {
		return 0, err
	}
	// Create a stub configuration so devices.configurationid FK is satisfied during enrollment.
	// This is needed because enrollment routes are decoupled from the legacy configurations table,
	// but the devices table still has a FK constraint on configurationid.
	_, _ = r.db.ExecContext(ctx, `
		INSERT INTO configurations (id, name, customerid, mainappid, permissive, type, settingsjson)
		VALUES ($1, $2, $3, $4, true, $5, '{}'::jsonb)
		ON CONFLICT (id) DO NOTHING`,
		id, strings.TrimSpace(req.Name)+" (enrollment)", customerID, resolved.VersionID, typ)
	// Add the bootstrap app to the stub configuration
	_, _ = r.db.ExecContext(ctx, `
		INSERT INTO configurationapplications (configurationid, applicationid, applicationversionid)
		SELECT $1, $2, $3
		WHERE NOT EXISTS (
			SELECT 1 FROM configurationapplications
			WHERE configurationid = $1 AND applicationid = $2
		)`, id, req.BootstrapApplicationID, resolved.VersionID)
	return id, nil
}

func (r *RouteRepository) Update(ctx context.Context, customerID, id int, req domain.UpdateRequest, resolved *domain.ResolvedBootstrap, containerAck *bool) error {
	cur, err := r.GetViewByID(ctx, customerID, id)
	if err != nil || cur == nil {
		return ErrNotFound
	}
	name := cur.Name
	if req.Name != nil {
		name = strings.TrimSpace(*req.Name)
	}
	desc := cur.Description
	if req.Description != nil {
		desc = *req.Description
	}
	treeID := cur.TargetNodeID
	if req.TargetNodeID != nil {
		treeID = *req.TargetNodeID
	} else if req.DefaultTreeNodeID != nil {
		treeID = *req.DefaultTreeNodeID
	}
	mode := cur.DeviceIdentityMode
	if req.DeviceIdentityMode != nil {
		mode = strings.TrimSpace(*req.DeviceIdentityMode)
	} else if req.DefaultDeviceIDMode != nil {
		mode = strings.TrimSpace(*req.DefaultDeviceIDMode)
	}
	if mode == "" {
		mode = "imei"
	}
	intent := cur.BootstrapIntent
	if req.BootstrapIntent != nil {
		intent = *req.BootstrapIntent
	}
	appID := cur.BootstrapApplicationID
	if req.BootstrapApplicationID != nil {
		appID = *req.BootstrapApplicationID
	}
	bootVer := cur.BootstrapVersionID
	if req.BootstrapVersionID != nil {
		bootVer = req.BootstrapVersionID
	}
	mainApp := cur.ResolvedMainAppVersionID
	if resolved != nil {
		v := resolved.VersionID
		mainApp = &v
	}
	ackAt := cur.ContainerPlacementAckAt
	if containerAck != nil {
		if *containerAck {
			now := time.Now()
			ackAt = &now
		} else {
			ackAt = nil
		}
	}
	wifiSSID := cur.WifiSSID
	if req.WifiSSID != nil {
		wifiSSID = *req.WifiSSID
	}
	wifiPassword := cur.WifiPassword
	if req.WifiPassword != nil {
		wifiPassword = *req.WifiPassword
	}
	wifiSecurityType := cur.WifiSecurityType
	if req.WifiSecurityType != nil {
		wifiSecurityType = *req.WifiSecurityType
	}
	qrParameters := cur.QRParameters
	if req.QRParameters != nil {
		qrParameters = *req.QRParameters
	}
	adminExtras := cur.AdminExtras
	if req.AdminExtras != nil {
		adminExtras = *req.AdminExtras
	}
	mobileEnrollment := cur.MobileEnrollment
	if req.MobileEnrollment != nil {
		mobileEnrollment = *req.MobileEnrollment
	}
	encryptDevice := cur.EncryptDevice
	if req.EncryptDevice != nil {
		encryptDevice = *req.EncryptDevice
	}
	res, err := r.db.ExecContext(ctx, `
		UPDATE enrollment_routes SET
			name=$1, description=$2, mainappid=$3,
			default_tree_node_id=$4, default_device_id_mode=$5,
			bootstrap_intent=$6, bootstrap_application_id=$7, bootstrap_version_id=$8,
			container_placement_ack_at=$9,
			wifi_ssid=$10, wifi_password=$11, wifi_security_type=$12,
			qr_parameters=$13, admin_extras=$14,
			mobile_enrollment=$15, encrypt_device=$16
		WHERE id=$17 AND customerid=$18`,
		name, desc, nullInt(mainApp), treeID, mode,
		intent, appID, nullIntPtr(bootVer), nullTime(ackAt),
		nullString(wifiSSID), nullString(wifiPassword), nullString(wifiSecurityType),
		nullString(qrParameters), nullString(adminExtras),
		mobileEnrollment, encryptDevice,
		id, customerID,
	)
	if err != nil {
		if strings.Contains(err.Error(), "enrollment_routes_name_customer_uidx") {
			return ErrDuplicateName
		}
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *RouteRepository) Delete(ctx context.Context, customerID, id int) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM enrollment_routes WHERE id = $1 AND customerid = $2`, id, customerID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *RouteRepository) DeleteImpact(ctx context.Context, customerID, routeID int) (*domain.EnrollmentDeleteImpact, error) {
	var out domain.EnrollmentDeleteImpact
	_ = r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM devices
		WHERE customerid = $1 AND enrollment_route_id = $2`, customerID, routeID).Scan(&out.HistoricalEnrolledCount)
	_ = r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM devices
		WHERE customerid = $1 AND enrollment_route_id = $2
		  AND enrolltime > NOW() - INTERVAL '24 hours'`, customerID, routeID).Scan(&out.EnrollingNowCount)
	agg := strconv.Itoa(routeID)
	_ = r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM domain_events
		WHERE event_type = 'enrollment_route.qr_viewed'
		  AND aggregate_id = $1
		  AND created_at > NOW() - INTERVAL '7 days'`, agg).Scan(&out.ActiveQrScans7d)
	return &out, nil
}

func (r *RouteRepository) TreeNodeBelongsToCustomer(ctx context.Context, customerID, nodeID int) (bool, error) {
	var ok int
	err := r.db.QueryRowContext(ctx, `
		SELECT 1 FROM device_tree_nodes WHERE id = $1 AND customerid = $2`, nodeID, customerID).Scan(&ok)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return err == nil, err
}

func (r *RouteRepository) NodePlacementKind(ctx context.Context, customerID, nodeID int) (string, error) {
	var child int
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM device_tree_nodes
		WHERE customerid = $1 AND parent_id = $2`, customerID, nodeID).Scan(&child)
	if err != nil {
		return "", err
	}
	if child > 0 {
		return "inheritable", nil
	}
	return "locked", nil
}

func (r *RouteRepository) ListTreeNodeOptions(ctx context.Context, customerID, heavyThreshold int) ([]domain.TreeNodeOption, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT n.id, n.name, COALESCE(n.path, ''), n.parent_id,
		       (SELECT COUNT(*)::int FROM devices d WHERE d.customerid = n.customerid AND d.tree_node_id = n.id) AS device_count,
		       (SELECT COUNT(*)::int FROM device_tree_nodes c WHERE c.parent_id = n.id AND c.customerid = n.customerid) AS child_count
		FROM device_tree_nodes n
		WHERE n.customerid = $1
		ORDER BY n.path`, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.TreeNodeOption
	for rows.Next() {
		var o domain.TreeNodeOption
		var parent sql.NullInt64
		var childCount int
		if err := rows.Scan(&o.ID, &o.Name, &o.Path, &parent, &o.DeviceCount, &childCount); err != nil {
			return nil, err
		}
		if parent.Valid {
			p := int(parent.Int64)
			o.ParentID = &p
		}
		if childCount > 0 {
			o.PlacementKind = "inheritable"
		} else {
			o.PlacementKind = "locked"
		}
		o.HeavilyLoaded = o.DeviceCount >= heavyThreshold
		out = append(out, o)
	}
	return out, rows.Err()
}

func (r *RouteRepository) ListBootstrapApps(ctx context.Context, customerID int) ([]domain.BootstrapAppOption, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT a.id, a.name, a.pkg
		FROM applications a
		WHERE a.customerid IS NULL OR a.customerid = $1
		ORDER BY lower(a.name)`, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	apps := make(map[int]*domain.BootstrapAppOption)
	var order []int
	for rows.Next() {
		var id int
		var name, pkg string
		if err := rows.Scan(&id, &name, &pkg); err != nil {
			return nil, err
		}
		apps[id] = &domain.BootstrapAppOption{ApplicationID: id, Name: name, Package: pkg}
		order = append(order, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	for _, appID := range order {
		vers, err := r.versionsForApp(ctx, appID)
		if err != nil {
			return nil, err
		}
		apps[appID].Versions = vers
	}
	out := make([]domain.BootstrapAppOption, 0, len(order))
	for _, id := range order {
		out = append(out, *apps[id])
	}
	return out, nil
}

func (r *RouteRepository) versionsForApp(ctx context.Context, applicationID int) ([]domain.BootstrapAppVersionOption, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT av.id, COALESCE(av.version, ''), av.versioncode, av.is_recommended
		FROM applicationversions av
		WHERE av.applicationid = $1
		ORDER BY av.versioncode DESC`, applicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var vers []domain.BootstrapAppVersionOption
	var maxCode int
	for rows.Next() {
		var v domain.BootstrapAppVersionOption
		if err := rows.Scan(&v.VersionID, &v.Version, &v.VersionCode, &v.IsRecommended); err != nil {
			return nil, err
		}
		if v.VersionCode > maxCode {
			maxCode = v.VersionCode
		}
		vers = append(vers, v)
	}
	for i := range vers {
		if vers[i].VersionCode == maxCode {
			vers[i].IsLatest = true
		}
	}
	return vers, rows.Err()
}

func (r *RouteRepository) RecordQRViewed(ctx context.Context, routeID int) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO domain_events (event_type, aggregate_id, payload)
		VALUES ('enrollment_route.qr_viewed', $1, '{}')`, strconv.Itoa(routeID))
	return err
}

func deviceMode(primary, fallback *string) string {
	if primary != nil && strings.TrimSpace(*primary) != "" {
		return strings.TrimSpace(*primary)
	}
	if fallback != nil && strings.TrimSpace(*fallback) != "" {
		return strings.TrimSpace(*fallback)
	}
	return "imei"
}

func nullStr(s *string) any {
	if s == nil {
		return nil
	}
	return *s
}

func nullInt(n *int) any {
	if n == nil || *n <= 0 {
		return nil
	}
	return *n
}

func nullIntPtr(n *int) any {
	if n == nil || *n <= 0 {
		return nil
	}
	return *n
}

func nullTime(t *time.Time) any {
	if t == nil {
		return nil
	}
	return *t
}

func nullBoolPtr(b *bool) any {
	if b == nil {
		return nil
	}
	return *b
}

func nullString(s string) any {
	if s == "" {
		return nil
	}
	return s
}
