package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/ZanDattSu/star-factory/auth/internal/model"
)

func (u *userRepository) Get(ctx context.Context, filter model.Filter) (*model.User, error) {
	var result *model.User

	err := WithTx(ctx, u.pool, func(tx pgx.Tx) error {
		user, err := u.getUser(ctx, tx, filter)
		if err != nil {
			return err
		}

		methods, err := u.getUserProviderMethods(ctx, tx, user.UUID)
		if err != nil {
			return err
		}

		user.Info.NotificationMethods = methods
		result = user

		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *userRepository) getUser(ctx context.Context, tx pgx.Tx, filter model.Filter) (*model.User, error) {
	if filter.UserUUID == nil && filter.UserLogin == nil {
		return nil, model.ErrInvalidUserFilter
	}

	var query string
	var args []any

	if filter.UserUUID != nil {
		query = `
			SELECT uuid, login, email, password, created_at, updated_at
			FROM users
			WHERE uuid = $1
		`
		id, err := uuid.Parse(*filter.UserUUID)
		if err != nil {
			return nil, model.ErrInvalidUserUUID
		}
		args = append(args, id)
	}

	if filter.UserLogin != nil {
		query = `
			SELECT uuid, login, email, password, created_at, updated_at
			FROM users
			WHERE login = $1
			LIMIT 1
		`
		args = append(args, *filter.UserLogin)
	}

	var (
		user      model.User
		createdAt pgtype.Timestamptz
		updatedAt pgtype.Timestamptz
	)

	err := tx.QueryRow(ctx, query, args...).Scan(
		&user.UUID,
		&user.Info.Login,
		&user.Info.Email,
		&user.Password,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	if createdAt.Valid {
		user.CreatedAt = &createdAt.Time
	}

	if updatedAt.Valid {
		user.UpdatedAt = &updatedAt.Time
	}

	return &user, nil
}

func (u *userRepository) getUserProviderMethods(ctx context.Context, tx pgx.Tx, userUUID uuid.UUID) ([]*model.NotificationMethod, error) {
	const query = `
		SELECT provider_name, target
		FROM notification_methods
		WHERE user_uuid = $1
	`

	rows, err := tx.Query(ctx, query, userUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to query notification methods: %w", err)
	}
	defer rows.Close()

	methods := make([]*model.NotificationMethod, 0)

	for rows.Next() {
		var m model.NotificationMethod
		err := rows.Scan(&m.ProviderName, &m.Target)
		if err != nil {
			return nil, fmt.Errorf("failed to scan method: %w", err)
		}
		methods = append(methods, &m)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows error: %w", rows.Err())
	}

	return methods, nil
}
