package user

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/ZanDattSu/star-factory/auth/internal/model"
	"github.com/ZanDattSu/star-factory/auth/internal/repository/converter"
)

func (u *userRepository) Create(ctx context.Context, user model.User) error {
	return WithTx(ctx, u.pool, func(tx pgx.Tx) error {
		repoUser := converter.UserToRepoModel(user)

		const insertUserSQL = `
			INSERT INTO users (
				uuid,
				login,
				email,
				password,
				created_at,
				updated_at
			)
			VALUES ($1, $2, $3, $4, $5, $6)
		`

		_, err := tx.Exec(ctx, insertUserSQL,
			repoUser.UUID,
			repoUser.Info.Login,
			repoUser.Info.Email,
			repoUser.Password,
			repoUser.CreatedAt,
			repoUser.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to insert user: %w", err)
		}

		// Если нет методов — выходим
		if repoUser.Info.NotificationMethods == nil {
			return nil
		}

		const insertMethodSQL = `
			INSERT INTO notification_methods (
				user_uuid,
				provider_name,
				target
			)
			VALUES ($1, $2, $3)
		`

		for _, m := range repoUser.Info.NotificationMethods {
			_, err := tx.Exec(ctx, insertMethodSQL,
				repoUser.UUID,
				m.ProviderName,
				m.Target,
			)
			if err != nil {
				return fmt.Errorf("failed to insert notification method: %w", err)
			}
		}

		return nil
	})
}
