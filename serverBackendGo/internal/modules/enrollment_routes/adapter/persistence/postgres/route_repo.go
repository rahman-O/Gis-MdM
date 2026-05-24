package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strings"

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
	ErrNotFound          = errors.New("enrollment route not found")
	ErrDuplicateName     = errors.New("duplicate enrollment route name")
	ErrInvalidBinding    = errors.New("invalid profile version or tree node")
)

func (r *RouteRepository) List(ctx context.Context, customerID int) ([]domain.Route, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT er.id, er.name, COALESCE(er.description, ''), COALESCE(er.qrcodekey, ''),
		       p.id, er.profile_version_id, pv.version_number,
		       er.default_tree_node_id, COALESCE(n.name, ''),
		       er.default_device_id_mode, er.mainappid
		FROM enrollment_routes er
		LEFT JOIN profile_versions pv ON pv.id = er.profile_version_id
		LEFT JOIN profiles p ON p.id = pv.profile_id
		LEFT JOIN device_tree_nodes n ON n.id = er.default_tree_node_id
		WHERE er.customerid = $1
		ORDER BY lower(er.name)`, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.Route
	for rows.Next() {
		var item domain.Route
		var profileID, pvID, treeID, verNum sql.NullInt64
		var mainApp sql.NullInt64
		if err := rows.Scan(
			&item.ID, &item.Name, &item.Description, &item.QRCodeKey,
			&profileID, &pvID, &verNum, &treeID, &item.DefaultTreeNodeName,
			&item.DefaultDeviceIDMode, &mainApp,
		); err != nil {
			return nil, err
		}
		if profileID.Valid {
			item.ProfileID = int(profileID.Int64)
		}
		if pvID.Valid {
			item.ProfileVersionID = int(pvID.Int64)
		}
		if verNum.Valid {
			v := int(verNum.Int64)
			item.ProfileVersionNumber = &v
		}
		if treeID.Valid {
			item.DefaultTreeNodeID = int(treeID.Int64)
		}
		if mainApp.Valid && mainApp.Int64 > 0 {
			m := int(mainApp.Int64)
			item.MainAppID = &m
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (r *RouteRepository) GetByID(ctx context.Context, customerID, id int) (*domain.RouteDetail, error) {
	var item domain.Route
	var profileID, pvID, treeID, verNum sql.NullInt64
	var mainApp sql.NullInt64
	err := r.db.QueryRowContext(ctx, `
		SELECT er.id, er.name, COALESCE(er.description, ''), COALESCE(er.qrcodekey, ''),
		       p.id, er.profile_version_id, pv.version_number,
		       er.default_tree_node_id, COALESCE(n.name, ''),
		       er.default_device_id_mode, er.mainappid, er.type
		FROM enrollment_routes er
		LEFT JOIN profile_versions pv ON pv.id = er.profile_version_id
		LEFT JOIN profiles p ON p.id = pv.profile_id
		LEFT JOIN device_tree_nodes n ON n.id = er.default_tree_node_id
		WHERE er.id = $1 AND er.customerid = $2`, id, customerID).Scan(
		&item.ID, &item.Name, &item.Description, &item.QRCodeKey,
		&profileID, &pvID, &verNum, &treeID, &item.DefaultTreeNodeName,
		&item.DefaultDeviceIDMode, &mainApp, &item.Type,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if profileID.Valid {
		item.ProfileID = int(profileID.Int64)
	}
	if pvID.Valid {
		item.ProfileVersionID = int(pvID.Int64)
	}
	if verNum.Valid {
		v := int(verNum.Int64)
		item.ProfileVersionNumber = &v
	}
	if treeID.Valid {
		item.DefaultTreeNodeID = int(treeID.Int64)
	}
	if mainApp.Valid && mainApp.Int64 > 0 {
		m := int(mainApp.Int64)
		item.MainAppID = &m
	}
	return &domain.RouteDetail{Route: item}, nil
}

func (r *RouteRepository) Create(ctx context.Context, customerID int, req domain.CreateRequest, qrcodeKey string) (int, error) {
	typ := 0
	if req.Type != nil {
		typ = *req.Type
	}
	mode := "imei"
	if req.DefaultDeviceIDMode != nil && strings.TrimSpace(*req.DefaultDeviceIDMode) != "" {
		mode = strings.TrimSpace(*req.DefaultDeviceIDMode)
	}
	var id int
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO enrollment_routes (
			customerid, name, description, qrcodekey, mainappid,
			profile_version_id, default_tree_node_id, default_device_id_mode, type
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING id`,
		customerID, strings.TrimSpace(req.Name), nullStr(req.Description), qrcodeKey,
		nullInt(req.MainAppID), nullIntPtr(req.ProfileVersionID), req.DefaultTreeNodeID, mode, typ,
	).Scan(&id)
	if err != nil && strings.Contains(err.Error(), "enrollment_routes_name_customer_uidx") {
		return 0, ErrDuplicateName
	}
	return id, err
}

func (r *RouteRepository) Update(ctx context.Context, customerID, id int, req domain.UpdateRequest, qrcodeKey *string) error {
	cur, err := r.GetByID(ctx, customerID, id)
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
	pvID := cur.ProfileVersionID
	if req.ProfileVersionID != nil {
		pvID = *req.ProfileVersionID
	}
	treeID := cur.DefaultTreeNodeID
	if req.DefaultTreeNodeID != nil {
		treeID = *req.DefaultTreeNodeID
	}
	mode := cur.DefaultDeviceIDMode
	if req.DefaultDeviceIDMode != nil {
		mode = strings.TrimSpace(*req.DefaultDeviceIDMode)
	}
	if mode == "" {
		mode = "imei"
	}
	mainApp := cur.MainAppID
	if req.MainAppID != nil {
		mainApp = req.MainAppID
	}
	qr := cur.QRCodeKey
	if qrcodeKey != nil {
		qr = *qrcodeKey
	}
	res, err := r.db.ExecContext(ctx, `
		UPDATE enrollment_routes SET
			name=$1, description=$2, qrcodekey=$3, mainappid=$4,
			profile_version_id=$5, default_tree_node_id=$6, default_device_id_mode=$7
		WHERE id=$8 AND customerid=$9`,
		name, desc, qr, nullInt(mainApp), pvID, treeID, mode, id, customerID,
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

func (r *RouteRepository) IsPublishedProfileVersion(ctx context.Context, customerID, profileVersionID int) (bool, error) {
	var ok int
	err := r.db.QueryRowContext(ctx, `
		SELECT 1 FROM profile_versions pv
		JOIN profiles p ON p.id = pv.profile_id
		WHERE pv.id = $1 AND p.customerid = $2 AND pv.status = 'published'`,
		profileVersionID, customerID,
	).Scan(&ok)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return err == nil, err
}

func (r *RouteRepository) ListPublishedProfileVersions(ctx context.Context, customerID int) ([]domain.PublishedProfileVersion, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT pv.id, p.id, p.name, pv.version_number, COALESCE(p.enabled, true), pv.mainappid
		FROM profile_versions pv
		JOIN profiles p ON p.id = pv.profile_id
		WHERE p.customerid = $1 AND pv.status = 'published'
		ORDER BY lower(p.name), pv.version_number DESC`, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.PublishedProfileVersion
	for rows.Next() {
		var item domain.PublishedProfileVersion
		var mainApp sql.NullInt64
		if err := rows.Scan(&item.ProfileVersionID, &item.ProfileID, &item.ProfileName, &item.VersionNumber, &item.ProfileEnabled, &mainApp); err != nil {
			return nil, err
		}
		if mainApp.Valid && mainApp.Int64 > 0 {
			m := int(mainApp.Int64)
			item.MainAppID = &m
		}
		out = append(out, item)
	}
	return out, rows.Err()
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

func nullStr(s *string) any {
	if s == nil {
		return nil
	}
	return *s
}

func nullInt(n *int) any {
	if n == nil || *n == 0 {
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
