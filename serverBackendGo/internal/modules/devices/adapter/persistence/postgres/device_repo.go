package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/gis-mdm/server-backend-go/internal/modules/devices/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/devices/port"
)

// DeviceRepository implements port.DeviceRepository.
type DeviceRepository struct {
	db *sql.DB
}

func NewDeviceRepository(db *sql.DB) *DeviceRepository {
	return &DeviceRepository{db: db}
}

var _ port.DeviceRepository = (*DeviceRepository)(nil)

func (r *DeviceRepository) LoadUserScope(ctx context.Context, userID int64) (*port.UserScope, error) {
	var scope port.UserScope
	scope.UserID = userID
	err := r.db.QueryRowContext(ctx, `
		SELECT customerid, alldevicesavailable FROM users WHERE id = $1`, userID).
		Scan(&scope.CustomerID, &scope.AllDevicesAvailable)
	if err != nil {
		return nil, err
	}
	return &scope, nil
}

const deviceAccessJoin = `
	LEFT JOIN devicegroups dg ON d.id = dg.deviceid
	LEFT JOIN groups g ON dg.groupid = g.id
	LEFT JOIN userdevicegroupsaccess access ON g.id = access.groupid AND access.userid = $2
`

const deviceAccessWhere = `
	d.customerid = $1
	AND ($3 = TRUE OR access.groupid IS NOT NULL)
`

func (r *DeviceRepository) Search(ctx context.Context, scope port.UserScope, req domain.SearchRequest) ([]domain.DeviceView, error) {
	args := []any{scope.CustomerID, scope.UserID, scope.AllDevicesAvailable}
	where := "WHERE " + deviceAccessWhere
	argN := 4
	searchFilters(req, &args, &where, &argN)
	offset := (req.PageNum - 1) * req.PageSize
	order := orderExpr(req)
	pageOrder := orderExprGrouped(req)
	pageSQL := fmt.Sprintf(`
		SELECT d.id
		FROM devices d
		%s
		%s
		%s
		GROUP BY d.id
		ORDER BY %s
		OFFSET $%d LIMIT $%d`, deviceAccessJoin, searchJoins(req), where, pageOrder, argN, argN+1)
	pageArgs := append(append([]any{}, args...), offset, req.PageSize)
	query := fmt.Sprintf(`
		SELECT d.id, d.number, d.description, d.lastupdate, d.configurationid, d.tree_node_id,
			d.enrollment_state, d.imei, d.phone,
			CASE
				WHEN (EXTRACT(EPOCH FROM NOW()) * 1000 - d.lastupdate) < (2 * 3600 * 1000) THEN 'green'
				WHEN (EXTRACT(EPOCH FROM NOW()) * 1000 - d.lastupdate) < (4 * 3600 * 1000) THEN 'yellow'
				ELSE 'red'
			END AS statuscode,
			d.infojson->>'model' AS model,
			d.infojson->>'androidVersion' AS android_version,
			d.infojson->>'serial' AS serial,
			d.infojson->>'launcherVersion' AS launcher_version,
			CASE WHEN d.infojson->>'batteryLevel' ~ '^\d+$' THEN (d.infojson->>'batteryLevel')::INT ELSE NULL END AS battery_level
		FROM devices d
		INNER JOIN (%s) page ON page.id = d.id
		ORDER BY %s`, pageSQL, order)
	args = pageArgs

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	var items []domain.DeviceView
	for rows.Next() {
		var v domain.DeviceView
		var desc, imei, phone, status sql.NullString
		var model, androidVersion, serial, launcherVersion sql.NullString
		var batteryLevel sql.NullInt64
		var lastUpdate sql.NullInt64
		var configID, treeNodeID sql.NullInt64
		var enrollmentState sql.NullString
		if err := rows.Scan(&v.ID, &v.Number, &desc, &lastUpdate, &configID, &treeNodeID, &enrollmentState, &imei, &phone, &status, &model, &androidVersion, &serial, &launcherVersion, &batteryLevel); err != nil {
			return nil, err
		}
		if desc.Valid {
			v.Description = &desc.String
		}
		if lastUpdate.Valid {
			v.LastUpdate = &lastUpdate.Int64
		}
		if configID.Valid {
			cid := int(configID.Int64)
			v.ConfigurationID = &cid
		}
		if treeNodeID.Valid {
			tid := int(treeNodeID.Int64)
			v.TreeNodeID = &tid
		}
		if enrollmentState.Valid {
			es := enrollmentState.String
			v.EnrollmentState = &es
		}
		if imei.Valid {
			v.IMEI = &imei.String
		}
		if phone.Valid {
			v.Phone = &phone.String
		}
		if status.Valid {
			v.StatusCode = &status.String
		}
		if model.Valid {
			v.Model = &model.String
		}
		if androidVersion.Valid {
			v.AndroidVersion = &androidVersion.String
		}
		if serial.Valid {
			v.Serial = &serial.String
		}
		if launcherVersion.Valid {
			v.LauncherVersion = &launcherVersion.String
		}
		if batteryLevel.Valid {
			bl := int(batteryLevel.Int64)
			v.BatteryLevel = &bl
		}
		v.Groups = []domain.LookupItem{}
		items = append(items, v)
		ids = append(ids, v.ID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return items, nil
	}
	groupMap, err := r.loadGroupsForDevices(ctx, ids)
	if err != nil {
		return nil, err
	}
	for i := range items {
		items[i].Groups = groupMap[items[i].ID]
	}
	return items, nil
}

func (r *DeviceRepository) loadGroupsForDevices(ctx context.Context, deviceIDs []int) (map[int][]domain.LookupItem, error) {
	out := make(map[int][]domain.LookupItem)
	if len(deviceIDs) == 0 {
		return out, nil
	}
	placeholders := make([]string, len(deviceIDs))
	args := make([]any, len(deviceIDs))
	for i, id := range deviceIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}
	q := fmt.Sprintf(`
		SELECT dg.deviceid, g.id, g.name
		FROM devicegroups dg
		INNER JOIN groups g ON dg.groupid = g.id
		WHERE dg.deviceid IN (%s)
		ORDER BY lower(g.name)`, strings.Join(placeholders, ","))
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var deviceID, groupID int
		var name string
		if err := rows.Scan(&deviceID, &groupID, &name); err != nil {
			return nil, err
		}
		n := name
		out[deviceID] = append(out[deviceID], domain.LookupItem{ID: groupID, Name: &n})
	}
	return out, rows.Err()
}

func (r *DeviceRepository) Count(ctx context.Context, scope port.UserScope, req domain.SearchRequest) (int64, error) {
	args := []any{scope.CustomerID, scope.UserID, scope.AllDevicesAvailable}
	where := "WHERE " + deviceAccessWhere
	argN := 4
	searchFilters(req, &args, &where, &argN)
	query := fmt.Sprintf(`
		SELECT COUNT(DISTINCT d.id)
		FROM devices d
		%s
		%s
		%s`, deviceAccessJoin, searchJoins(req), where)
	var n int64
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&n)
	return n, err
}

func (r *DeviceRepository) ListConfigurations(ctx context.Context, customerID int) (map[int]domain.ConfigurationView, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, permissive FROM configurations WHERE customerid = $1 ORDER BY lower(name)`, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make(map[int]domain.ConfigurationView)
	for rows.Next() {
		var cv domain.ConfigurationView
		var name string
		var permissive bool
		if err := rows.Scan(&cv.ID, &name, &permissive); err != nil {
			return nil, err
		}
		cv.Name = &name
		cv.PermissiveMode = &permissive
		out[cv.ID] = cv
	}
	return out, rows.Err()
}

func (r *DeviceRepository) GetByNumber(ctx context.Context, scope port.UserScope, number string) (*domain.DeviceView, error) {
	args := []any{scope.CustomerID, scope.UserID, scope.AllDevicesAvailable, number}
	query := fmt.Sprintf(`
		SELECT DISTINCT d.id, d.number, d.description, d.lastupdate, d.configurationid, d.imei, d.phone,
			d.custom1, d.custom2, d.custom3, d.oldnumber,
			d.info, d.infojson, d.enrolltime, d.publicip,
			CASE
				WHEN (EXTRACT(EPOCH FROM NOW()) * 1000 - d.lastupdate) < (2 * 3600 * 1000) THEN 'green'
				WHEN (EXTRACT(EPOCH FROM NOW()) * 1000 - d.lastupdate) < (4 * 3600 * 1000) THEN 'yellow'
				ELSE 'red'
			END
		FROM devices d
		%s
		WHERE d.customerid = $1 AND ($3 = TRUE OR access.groupid IS NOT NULL) AND lower(d.number) = lower($4)
		LIMIT 1`, deviceAccessJoin)
	var v domain.DeviceView
	var desc, imei, phone, status sql.NullString
	var custom1, custom2, custom3, oldNumber sql.NullString
	var info sql.NullString
	var infojson []byte
	var enrollTime sql.NullInt64
	var publicIP sql.NullString
	var lastUpdate sql.NullInt64
	var configID sql.NullInt64
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&v.ID, &v.Number, &desc, &lastUpdate, &configID, &imei, &phone,
		&custom1, &custom2, &custom3, &oldNumber,
		&info, &infojson, &enrollTime, &publicIP, &status)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	fillNulls(&v, desc, lastUpdate, configID, imei, phone, status)
	if custom1.Valid {
		v.Custom1 = &custom1.String
	}
	if custom2.Valid {
		v.Custom2 = &custom2.String
	}
	if custom3.Valid {
		v.Custom3 = &custom3.String
	}
	if oldNumber.Valid {
		v.OldNumber = &oldNumber.String
	}
	v.Info = parseDeviceInfo(info, infojson, enrollTime, publicIP)
	groups, _ := r.loadGroupsForDevices(ctx, []int{v.ID})
	v.Groups = groups[v.ID]
	return &v, nil
}

func (r *DeviceRepository) GetByID(ctx context.Context, customerID int, id int) (*domain.DeviceView, error) {
	var v domain.DeviceView
	var desc, imei, phone sql.NullString
	var lastUpdate sql.NullInt64
	var configID sql.NullInt64
	err := r.db.QueryRowContext(ctx, `
		SELECT id, number, description, lastupdate, configurationid, imei, phone
		FROM devices WHERE id = $1 AND customerid = $2`, id, customerID).
		Scan(&v.ID, &v.Number, &desc, &lastUpdate, &configID, &imei, &phone)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	fillNulls(&v, desc, lastUpdate, configID, imei, phone, sql.NullString{})
	groups, _ := r.loadGroupsForDevices(ctx, []int{v.ID})
	v.Groups = groups[v.ID]
	return &v, nil
}

func fillNulls(v *domain.DeviceView, desc sql.NullString, lastUpdate, configID sql.NullInt64, imei, phone, status sql.NullString) {
	if desc.Valid {
		v.Description = &desc.String
	}
	if lastUpdate.Valid {
		v.LastUpdate = &lastUpdate.Int64
	}
	if configID.Valid {
		cid := int(configID.Int64)
		v.ConfigurationID = &cid
	}
	if imei.Valid {
		v.IMEI = &imei.String
	}
	if phone.Valid {
		v.Phone = &phone.String
	}
	if status.Valid {
		v.StatusCode = &status.String
	}
	if v.Groups == nil {
		v.Groups = []domain.LookupItem{}
	}
}

func (r *DeviceRepository) ExistsNumber(ctx context.Context, customerID int, number string, excludeID int) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM devices WHERE customerid = $1 AND lower(number) = lower($2) AND ($3 = 0 OR id <> $3)
		)`, customerID, number, excludeID).Scan(&exists)
	return exists, err
}

func (r *DeviceRepository) CountDevices(ctx context.Context, customerID int) (int64, error) {
	var n int64
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM devices WHERE customerid = $1`, customerID).Scan(&n)
	return n, err
}

func (r *DeviceRepository) DeviceLimit(ctx context.Context, customerID int) (int, error) {
	var limit sql.NullInt64
	err := r.db.QueryRowContext(ctx, `SELECT devicelimit FROM customers WHERE id = $1`, customerID).Scan(&limit)
	if err != nil {
		return 0, err
	}
	if !limit.Valid {
		return 0, nil
	}
	return int(limit.Int64), nil
}

func (r *DeviceRepository) Insert(ctx context.Context, customerID int, d domain.SaveDevice) (int, error) {
	number := strings.TrimSpace(ptrStr(d.Number))
	var id int
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO devices (number, description, lastupdate, configurationid, customerid, enrolltime)
		VALUES ($1, $2, EXTRACT(EPOCH FROM NOW())::BIGINT * 1000, $3, $4, EXTRACT(EPOCH FROM NOW())::BIGINT * 1000)
		RETURNING id`, number, ptrStr(d.Description), *d.ConfigurationID, customerID).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, r.replaceDeviceGroups(ctx, id, d.Groups)
}

func (r *DeviceRepository) Update(ctx context.Context, customerID int, d domain.SaveDevice) error {
	if d.ID == nil {
		return fmt.Errorf("missing device id")
	}
	_, err := r.db.ExecContext(ctx, `
		UPDATE devices SET
			number = COALESCE($1, number),
			description = COALESCE($2, description),
			configurationid = COALESCE($3, configurationid),
			lastupdate = EXTRACT(EPOCH FROM NOW())::BIGINT * 1000
		WHERE id = $4 AND customerid = $5`,
		ptrStr(d.Number), ptrStr(d.Description), d.ConfigurationID, *d.ID, customerID)
	if err != nil {
		return err
	}
	return r.replaceDeviceGroups(ctx, *d.ID, d.Groups)
}

func (r *DeviceRepository) replaceDeviceGroups(ctx context.Context, deviceID int, groups []domain.LookupItem) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM devicegroups WHERE deviceid = $1`, deviceID)
	if err != nil {
		return err
	}
	for _, g := range groups {
		if g.ID <= 0 {
			continue
		}
		_, err := r.db.ExecContext(ctx, `
			INSERT INTO devicegroups (deviceid, groupid) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
			deviceID, g.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *DeviceRepository) UpdateConfigurationBulk(ctx context.Context, customerID int, ids []int, configID int) error {
	for _, id := range ids {
		_, err := r.db.ExecContext(ctx, `
			UPDATE devices SET configurationid = $1, lastupdate = EXTRACT(EPOCH FROM NOW())::BIGINT * 1000
			WHERE id = $2 AND customerid = $3`, configID, id, customerID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *DeviceRepository) Delete(ctx context.Context, customerID int, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM devices WHERE id = $1 AND customerid = $2`, id, customerID)
	return err
}

func (r *DeviceRepository) DeleteBulk(ctx context.Context, customerID int, ids []int) error {
	for _, id := range ids {
		if err := r.Delete(ctx, customerID, id); err != nil {
			return err
		}
	}
	return nil
}

func (r *DeviceRepository) UpdateGroupBulk(ctx context.Context, customerID int, req domain.GroupBulkRequest) error {
	for _, id := range req.IDs {
		var devCustomer int
		if err := r.db.QueryRowContext(ctx, `SELECT customerid FROM devices WHERE id = $1`, id).Scan(&devCustomer); err != nil {
			return err
		}
		if devCustomer != customerID {
			continue
		}
		if strings.EqualFold(req.Action, "set") {
			_ = r.replaceDeviceGroups(ctx, id, req.Groups)
		} else {
			_, _ = r.db.ExecContext(ctx, `DELETE FROM devicegroups WHERE deviceid = $1`, id)
		}
	}
	return nil
}

func (r *DeviceRepository) Autocomplete(ctx context.Context, scope port.UserScope, filter string, limit int) ([]domain.LookupItem, error) {
	req := domain.SearchRequest{PageNum: 1, PageSize: limit, Value: &filter}
	items, err := r.Search(ctx, scope, req)
	if err != nil {
		return nil, err
	}
	out := make([]domain.LookupItem, 0, len(items))
	for _, d := range items {
		num := d.Number
		out = append(out, domain.LookupItem{ID: d.ID, Name: &num})
	}
	return out, nil
}

func (r *DeviceRepository) UpdateDescription(ctx context.Context, customerID int, id int, description string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE devices SET description = $1, lastupdate = EXTRACT(EPOCH FROM NOW())::BIGINT * 1000
		WHERE id = $2 AND customerid = $3`, description, id, customerID)
	return err
}

func (r *DeviceRepository) ListAppSettings(ctx context.Context, deviceID int) ([]domain.AppSetting, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT applicationpkg, name, type, value FROM deviceapplicationsettings WHERE deviceid = $1`, deviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.AppSetting
	for rows.Next() {
		var s domain.AppSetting
		var pkg, name, typ, val sql.NullString
		if err := rows.Scan(&pkg, &name, &typ, &val); err != nil {
			return nil, err
		}
		if pkg.Valid {
			s.ApplicationPkg = &pkg.String
		}
		if name.Valid {
			s.Name = &name.String
		}
		if typ.Valid {
			s.Type = &typ.String
		}
		if val.Valid {
			s.Value = &val.String
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func (r *DeviceRepository) SaveAppSettings(ctx context.Context, deviceID int, settings []domain.AppSetting) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM deviceapplicationsettings WHERE deviceid = $1`, deviceID)
	if err != nil {
		return err
	}
	for _, s := range settings {
		pkg := ptrStr(s.ApplicationPkg)
		if pkg == "" {
			continue
		}
		_, err := r.db.ExecContext(ctx, `
			INSERT INTO deviceapplicationsettings (deviceid, applicationpkg, name, type, value)
			VALUES ($1, $2, $3, $4, $5)`,
			deviceID, pkg, ptrStr(s.Name), ptrStr(s.Type), ptrStr(s.Value))
		if err != nil {
			return err
		}
	}
	return nil
}

func ptrStr(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func (r *DeviceRepository) MoveTreeNode(ctx context.Context, customerID int, deviceID int, treeNodeID int) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE devices d SET tree_node_id = $1
		FROM device_tree_nodes n
		WHERE d.id = $2 AND d.customerid = $3
		  AND n.id = $1 AND n.customerid = $3`,
		treeNodeID, deviceID, customerID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
