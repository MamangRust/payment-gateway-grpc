// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: role.sql

package db

import (
	"context"
	"database/sql"
)

const createRole = `-- name: CreateRole :one
INSERT INTO roles (
    role_name, 
    created_at, 
    updated_at
) VALUES (
    $1, 
    current_timestamp, 
    current_timestamp
) RETURNING 
    role_id, 
    role_name, 
    created_at, 
    updated_at, 
    deleted_at
`

func (q *Queries) CreateRole(ctx context.Context, roleName string) (*Role, error) {
	row := q.db.QueryRowContext(ctx, createRole, roleName)
	var i Role
	err := row.Scan(
		&i.RoleID,
		&i.RoleName,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return &i, err
}

const deletePermanentRole = `-- name: DeletePermanentRole :exec
DELETE FROM roles
WHERE 
    role_id = $1
`

func (q *Queries) DeletePermanentRole(ctx context.Context, roleID int32) error {
	_, err := q.db.ExecContext(ctx, deletePermanentRole, roleID)
	return err
}

const getActiveRoles = `-- name: GetActiveRoles :many
SELECT 
    role_id, 
    role_name, 
    created_at, 
    updated_at, 
    deleted_at 
FROM 
    roles
WHERE 
    deleted_at IS NULL
    AND role_name ILIKE '%' || $1 || '%'
ORDER BY 
    created_at ASC
LIMIT $2 OFFSET $3
`

type GetActiveRolesParams struct {
	Column1 sql.NullString `json:"column_1"`
	Limit   int32          `json:"limit"`
	Offset  int32          `json:"offset"`
}

func (q *Queries) GetActiveRoles(ctx context.Context, arg GetActiveRolesParams) ([]*Role, error) {
	rows, err := q.db.QueryContext(ctx, getActiveRoles, arg.Column1, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*Role
	for rows.Next() {
		var i Role
		if err := rows.Scan(
			&i.RoleID,
			&i.RoleName,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRole = `-- name: GetRole :one
SELECT 
    role_id, 
    role_name, 
    created_at, 
    updated_at, 
    deleted_at 
FROM 
    roles
WHERE 
    role_id = $1
`

func (q *Queries) GetRole(ctx context.Context, roleID int32) (*Role, error) {
	row := q.db.QueryRowContext(ctx, getRole, roleID)
	var i Role
	err := row.Scan(
		&i.RoleID,
		&i.RoleName,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return &i, err
}

const getRoleByName = `-- name: GetRoleByName :one
SELECT 
    role_id, 
    role_name, 
    created_at, 
    updated_at, 
    deleted_at 
FROM 
    roles
WHERE 
    role_name = $1
`

func (q *Queries) GetRoleByName(ctx context.Context, roleName string) (*Role, error) {
	row := q.db.QueryRowContext(ctx, getRoleByName, roleName)
	var i Role
	err := row.Scan(
		&i.RoleID,
		&i.RoleName,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return &i, err
}

const getRoles = `-- name: GetRoles :many
SELECT 
    role_id, 
    role_name, 
    created_at, 
    updated_at, 
    deleted_at 
FROM 
    roles
WHERE 
    role_name ILIKE '%' || $1 || '%'
ORDER BY 
    created_at ASC
LIMIT $2 OFFSET $3
`

type GetRolesParams struct {
	Column1 sql.NullString `json:"column_1"`
	Limit   int32          `json:"limit"`
	Offset  int32          `json:"offset"`
}

func (q *Queries) GetRoles(ctx context.Context, arg GetRolesParams) ([]*Role, error) {
	rows, err := q.db.QueryContext(ctx, getRoles, arg.Column1, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*Role
	for rows.Next() {
		var i Role
		if err := rows.Scan(
			&i.RoleID,
			&i.RoleName,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTrashedRoles = `-- name: GetTrashedRoles :many
SELECT 
    role_id, 
    role_name, 
    created_at, 
    updated_at, 
    deleted_at 
FROM 
    roles
WHERE 
    deleted_at IS NOT NULL
    AND role_name ILIKE '%' || $1 || '%'
ORDER BY 
    deleted_at DESC
LIMIT $2 OFFSET $3
`

type GetTrashedRolesParams struct {
	Column1 sql.NullString `json:"column_1"`
	Limit   int32          `json:"limit"`
	Offset  int32          `json:"offset"`
}

func (q *Queries) GetTrashedRoles(ctx context.Context, arg GetTrashedRolesParams) ([]*Role, error) {
	rows, err := q.db.QueryContext(ctx, getTrashedRoles, arg.Column1, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*Role
	for rows.Next() {
		var i Role
		if err := rows.Scan(
			&i.RoleID,
			&i.RoleName,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserRoles = `-- name: GetUserRoles :many
SELECT 
    r.role_id, 
    r.role_name, 
    r.created_at, 
    r.updated_at, 
    r.deleted_at 
FROM 
    roles r
JOIN 
    user_roles ur ON ur.role_id = r.role_id
WHERE 
    ur.user_id = $1
ORDER BY 
    r.created_at ASC
`

func (q *Queries) GetUserRoles(ctx context.Context, userID int32) ([]*Role, error) {
	rows, err := q.db.QueryContext(ctx, getUserRoles, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*Role
	for rows.Next() {
		var i Role
		if err := rows.Scan(
			&i.RoleID,
			&i.RoleName,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const restoreRole = `-- name: RestoreRole :exec
UPDATE roles
SET 
    deleted_at = NULL
WHERE 
    role_id = $1
`

func (q *Queries) RestoreRole(ctx context.Context, roleID int32) error {
	_, err := q.db.ExecContext(ctx, restoreRole, roleID)
	return err
}

const trashRole = `-- name: TrashRole :exec
UPDATE roles
SET 
    deleted_at = current_timestamp
WHERE 
    role_id = $1
`

func (q *Queries) TrashRole(ctx context.Context, roleID int32) error {
	_, err := q.db.ExecContext(ctx, trashRole, roleID)
	return err
}

const updateRole = `-- name: UpdateRole :one
UPDATE roles
SET 
    role_name = $2,
    updated_at = current_timestamp
WHERE 
    role_id = $1
RETURNING 
    role_id, 
    role_name, 
    created_at, 
    updated_at, 
    deleted_at
`

type UpdateRoleParams struct {
	RoleID   int32  `json:"role_id"`
	RoleName string `json:"role_name"`
}

func (q *Queries) UpdateRole(ctx context.Context, arg UpdateRoleParams) (*Role, error) {
	row := q.db.QueryRowContext(ctx, updateRole, arg.RoleID, arg.RoleName)
	var i Role
	err := row.Scan(
		&i.RoleID,
		&i.RoleName,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return &i, err
}
