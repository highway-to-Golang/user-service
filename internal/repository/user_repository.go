package repository

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/highway-to-Golang/user-service/internal/database"
	"github.com/highway-to-Golang/user-service/internal/domain"
	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	db   *database.DB
	goqu *goqu.Database
}

func NewUserRepository(db *database.DB) *UserRepository {
	goquDB := goqu.New("postgres", db.SQLDB())

	return &UserRepository{
		db:   db,
		goqu: goquDB,
	}
}

func (r *UserRepository) Create(ctx context.Context, user domain.User) error {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	query, args, err := r.goqu.Insert("users").
		Cols("id", "name", "email", "role", "created_at", "updated_at").
		Vals(goqu.Vals{user.ID, user.Name, user.Email, user.Role, user.CreatedAt, user.UpdatedAt}).
		ToSQL()

	if err != nil {
		slog.Error("failed to build insert query", "error", err)
		return fmt.Errorf("failed to build insert query: %w", err)
	}

	slog.Debug("executing insert query", "query", query, "args", args)

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		slog.Error("failed to create user", "error", err, "user_id", user.ID)
		return fmt.Errorf("failed to create user: %w", err)
	}

	slog.Info("user created successfully", "user_id", user.ID, "email", user.Email, "created_at", user.CreatedAt)
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (domain.User, error) {
	query, args, err := r.goqu.From("users").
		Select("id", "name", "email", "role", "created_at", "updated_at").
		Where(goqu.C("id").Eq(id)).
		ToSQL()

	if err != nil {
		slog.Error("failed to build select query", "error", err)
		return domain.User{}, fmt.Errorf("failed to build select query: %w", err)
	}

	slog.Debug("executing select query", "query", query, "args", args)

	var user domain.User
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			slog.Warn("user not found", "user_id", id)
			return domain.User{}, domain.ErrNotFound
		}
		slog.Error("failed to get user", "error", err, "user_id", id)
		return domain.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	slog.Info("user retrieved successfully", "user_id", id)
	return user, nil
}

func (r *UserRepository) GetAll(ctx context.Context) ([]domain.User, error) {
	query, args, err := r.goqu.From("users").
		Select("id", "name", "email", "role", "created_at", "updated_at").
		Order(goqu.C("created_at").Desc()).
		ToSQL()

	if err != nil {
		slog.Error("failed to build select all query", "error", err)
		return nil, fmt.Errorf("failed to build select all query: %w", err)
	}

	slog.Debug("executing select all query", "query", query, "args", args)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		slog.Error("failed to get users", "error", err)
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			slog.Error("failed to scan user", "error", err)
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		slog.Error("error during rows iteration", "error", err)
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	slog.Info("users retrieved successfully", "count", len(users))
	return users, nil
}

func (r *UserRepository) Update(ctx context.Context, id string, user domain.User) error {
	user.UpdatedAt = time.Now()

	query, args, err := r.goqu.Update("users").
		Set(goqu.Record{
			"name":       user.Name,
			"email":      user.Email,
			"role":       user.Role,
			"updated_at": user.UpdatedAt,
		}).
		Where(goqu.C("id").Eq(id)).
		ToSQL()

	if err != nil {
		slog.Error("failed to build update query", "error", err)
		return fmt.Errorf("failed to build update query: %w", err)
	}

	slog.Debug("executing update query", "query", query, "args", args)

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		slog.Error("failed to update user", "error", err, "user_id", id)
		return fmt.Errorf("failed to update user: %w", err)
	}

	if result.RowsAffected() == 0 {
		slog.Warn("user not found for update", "user_id", id)
		return domain.ErrNotFound
	}

	slog.Info("user updated successfully", "user_id", id, "updated_at", user.UpdatedAt)
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	query, args, err := r.goqu.Delete("users").
		Where(goqu.C("id").Eq(id)).
		ToSQL()

	if err != nil {
		slog.Error("failed to build delete query", "error", err)
		return fmt.Errorf("failed to build delete query: %w", err)
	}

	slog.Debug("executing delete query", "query", query, "args", args)

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		slog.Error("failed to delete user", "error", err, "user_id", id)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		slog.Warn("user not found for deletion", "user_id", id)
		return domain.ErrNotFound
	}

	slog.Info("user deleted successfully", "user_id", id)
	return nil
}
