package migration

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func runSchemaMigrations(ctx context.Context, db *pgxpool.Pool) error {
	_, err := db.Exec(
		ctx,
		`
			CREATE TABLE IF NOT EXISTS schema_migrations (
				version TEXT PRIMARY KEY,
				name TEXT NOT NULL,
				applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
			);
		`,
	)

	return err
}

func getAppliedMigrations(ctx context.Context, db *pgxpool.Pool) (map[string]bool, error) {
	rows, err := db.Query(
		ctx,
		`
			SELECT
				version
			FROM
				schema_migrations
		`,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	applied := map[string]bool{}

	for rows.Next() {
		var version string

		if err := rows.Scan(&version); err != nil {
			return nil, err
		}

		applied[version] = true
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return applied, nil
}

func applyMigration(ctx context.Context, db *pgxpool.Pool, migration Migration) error {
	sqlContent, err := os.ReadFile(migration.UpFile)
	if err != nil {
		return err
	}

	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, string(sqlContent)); err != nil {
		return err
	}

	if _, err := tx.Exec(
		ctx,
		`
		 	INSERT INTO schema_migrations (version, name)
			VALUES ($1, $2)
		`,
		migration.Version,
		migration.Name,
	); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func rollbackMigration(ctx context.Context, db *pgxpool.Pool, migration Migration) error {
	sqlContent, err := os.ReadFile(migration.DownFile)
	if err != nil {
		return err
	}

	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, string(sqlContent)); err != nil {
		return fmt.Errorf("erro ao executar down da migration %s_%s: %w", migration.Version, migration.Name, err)
	}

	if _, err := tx.Exec(
		ctx,
		`
			DELETE FROM schema_migrations
			WHERE version = $1	
		`,
		migration.Version,
	); err != nil {
		return fmt.Errorf("erro ao remover registro da migration %s_%s: %w", migration.Version, migration.Name, err)
	}

	return tx.Commit(ctx)
}
