package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/gis-mdm/server-backend-go/internal/modules/device_tree/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/device_tree/port"
)

// TreeRepository implements port.TreeRepository.
type TreeRepository struct {
	db *sql.DB
}

func NewTreeRepository(db *sql.DB) *TreeRepository {
	return &TreeRepository{db: db}
}

var _ port.TreeRepository = (*TreeRepository)(nil)

func (r *TreeRepository) EnsureRoot(ctx context.Context, customerID int) (int, error) {
	var id int
	err := r.db.QueryRowContext(ctx, `
		SELECT id FROM device_tree_nodes
		WHERE customerid = $1 AND parent_id IS NULL
		ORDER BY id LIMIT 1`, customerID).Scan(&id)
	if err == nil {
		return id, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}
	err = r.db.QueryRowContext(ctx, `
		INSERT INTO device_tree_nodes (customerid, parent_id, name, sort_order, path, depth)
		VALUES ($1, NULL, $2, 0, '', 0)
		RETURNING id`, customerID, domain.RootFolderName).Scan(&id)
	if err != nil {
		return 0, err
	}
	path := fmt.Sprintf("/%d/", id)
	_, err = r.db.ExecContext(ctx, `
		UPDATE device_tree_nodes SET path = $1, depth = 0 WHERE id = $2`, path, id)
	return id, err
}

func (r *TreeRepository) ListByCustomer(ctx context.Context, customerID int) ([]domain.TreeNode, error) {
	if _, err := r.EnsureRoot(ctx, customerID); err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, `
		SELECT n.id, n.parent_id, n.name, n.sort_order, n.path, n.depth,
			(SELECT COUNT(*)::INT FROM devices d
			 WHERE d.customerid = n.customerid AND d.tree_node_id = n.id) AS device_count
		FROM device_tree_nodes n
		WHERE n.customerid = $1
		ORDER BY n.depth, n.sort_order, lower(n.name)`, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.TreeNode
	for rows.Next() {
		var n domain.TreeNode
		var parent sql.NullInt64
		if err := rows.Scan(&n.ID, &parent, &n.Name, &n.SortOrder, &n.Path, &n.Depth, &n.DeviceCount); err != nil {
			return nil, err
		}
		if parent.Valid {
			pid := int(parent.Int64)
			n.ParentID = &pid
		}
		out = append(out, n)
	}
	return out, rows.Err()
}

func (r *TreeRepository) GetByID(ctx context.Context, customerID, id int) (*domain.TreeNode, error) {
	var n domain.TreeNode
	var parent sql.NullInt64
	err := r.db.QueryRowContext(ctx, `
		SELECT n.id, n.parent_id, n.name, n.sort_order, n.path, n.depth,
			(SELECT COUNT(*)::INT FROM devices d
			 WHERE d.customerid = n.customerid AND d.tree_node_id = n.id)
		FROM device_tree_nodes n
		WHERE n.customerid = $1 AND n.id = $2`, customerID, id).
		Scan(&n.ID, &parent, &n.Name, &n.SortOrder, &n.Path, &n.Depth, &n.DeviceCount)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if parent.Valid {
		pid := int(parent.Int64)
		n.ParentID = &pid
	}
	return &n, nil
}

func (r *TreeRepository) Create(ctx context.Context, customerID int, req domain.CreateNodeRequest) (*domain.TreeNode, error) {
	parent, err := r.GetByID(ctx, customerID, req.ParentID)
	if err != nil {
		return nil, domain.ErrInvalidParent
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, domain.ErrDuplicateName
	}
	var id int
	err = r.db.QueryRowContext(ctx, `
		INSERT INTO device_tree_nodes (customerid, parent_id, name, sort_order, path, depth)
		VALUES ($1, $2, $3, $4, '', $5)
		RETURNING id`,
		customerID, req.ParentID, name, req.SortOrder, parent.Depth+1).Scan(&id)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, domain.ErrDuplicateName
		}
		return nil, err
	}
	path := parent.Path + fmt.Sprintf("%d/", id)
	_, err = r.db.ExecContext(ctx, `UPDATE device_tree_nodes SET path = $1 WHERE id = $2`, path, id)
	if err != nil {
		return nil, err
	}
	return r.GetByID(ctx, customerID, id)
}

func (r *TreeRepository) Update(ctx context.Context, customerID, id int, req domain.UpdateNodeRequest) (*domain.TreeNode, error) {
	node, err := r.GetByID(ctx, customerID, id)
	if err != nil {
		return nil, err
	}
	if req.ParentID != nil && *req.ParentID != 0 {
		if *req.ParentID == id {
			return nil, domain.ErrCycle
		}
		newParent, err := r.GetByID(ctx, customerID, *req.ParentID)
		if err != nil {
			return nil, domain.ErrInvalidParent
		}
		if strings.HasPrefix(newParent.Path, node.Path) {
			return nil, domain.ErrCycle
		}
		if err := r.reparentSubtree(ctx, customerID, node, newParent); err != nil {
			return nil, err
		}
	}
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, domain.ErrDuplicateName
		}
		_, err := r.db.ExecContext(ctx, `
			UPDATE device_tree_nodes SET name = $1 WHERE customerid = $2 AND id = $3`,
			name, customerID, id)
		if err != nil {
			if isUniqueViolation(err) {
				return nil, domain.ErrDuplicateName
			}
			return nil, err
		}
	}
	if req.SortOrder != nil {
		_, err := r.db.ExecContext(ctx, `
			UPDATE device_tree_nodes SET sort_order = $1 WHERE customerid = $2 AND id = $3`,
			*req.SortOrder, customerID, id)
		if err != nil {
			return nil, err
		}
	}
	return r.GetByID(ctx, customerID, id)
}

func (r *TreeRepository) reparentSubtree(ctx context.Context, customerID int, node, newParent *domain.TreeNode) error {
	oldPath := node.Path
	newPath := newParent.Path + fmt.Sprintf("%d/", node.ID)
	newDepth := newParent.Depth + 1
	depthDelta := newDepth - node.Depth

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		UPDATE device_tree_nodes
		SET parent_id = $1, path = $2, depth = $3
		WHERE customerid = $4 AND id = $5`,
		newParent.ID, newPath, newDepth, customerID, node.ID)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `
		UPDATE device_tree_nodes
		SET path = $1 || SUBSTRING(path FROM LENGTH($2) + 1),
		    depth = depth + $3
		WHERE customerid = $4 AND path LIKE $5 AND id <> $6`,
		newPath, oldPath, depthDelta, customerID, oldPath+"%", node.ID)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (r *TreeRepository) CountDevicesInSubtree(ctx context.Context, customerID, nodeID int) (int, error) {
	node, err := r.GetByID(ctx, customerID, nodeID)
	if err != nil {
		return 0, err
	}
	var count int
	err = r.db.QueryRowContext(ctx, `
		SELECT COUNT(*)::INT FROM devices d
		INNER JOIN device_tree_nodes n ON n.id = d.tree_node_id
		WHERE d.customerid = $1 AND n.path LIKE $2`,
		customerID, node.Path+"%").Scan(&count)
	return count, err
}

func (r *TreeRepository) DeleteWithRelocation(ctx context.Context, customerID, id, targetNodeID int) error {
	node, err := r.GetByID(ctx, customerID, id)
	if err != nil {
		return err
	}
	if node.ParentID == nil {
		return domain.ErrCannotDeleteRoot
	}
	target, err := r.GetByID(ctx, customerID, targetNodeID)
	if err != nil {
		return domain.ErrInvalidParent
	}
	if target.ID == id || strings.HasPrefix(target.Path, node.Path) {
		return domain.ErrInvalidParent
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		UPDATE devices d SET tree_node_id = $1
		FROM device_tree_nodes n
		WHERE d.tree_node_id = n.id
		  AND d.customerid = $2
		  AND n.customerid = $2
		  AND n.path LIKE $3`,
		targetNodeID, customerID, node.Path+"%")
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `
		DELETE FROM device_tree_nodes
		WHERE customerid = $1 AND path LIKE $2`,
		customerID, node.Path+"%")
	if err != nil {
		return err
	}
	return tx.Commit()
}

func isUniqueViolation(err error) bool {
	return err != nil && strings.Contains(err.Error(), "unique")
}
