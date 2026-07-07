package migration

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func runUp(ctx context.Context, db *pgxpool.Pool, migrations []Migration) error {
	applied, err := getAppliedMigrations(ctx, db)
	if err != nil {
		return err
	}

	for _, migration := range migrations {
		if applied[migration.Version] {
			continue
		}

		if err := applyMigration(ctx, db, migration); err != nil {
			return err
		}
	}

	return nil
}

func runDown(ctx context.Context, db *pgxpool.Pool, migrations []Migration) error {
	applied, err := getAppliedMigrations(ctx, db)

	if err != nil {
		return err
	}

	for i := len(migrations) - 1; i >= 0; i-- {
		migration := migrations[i]

		if !applied[migration.Version] {
			continue
		}

		return rollbackMigration(ctx, db, migration)
	}

	return nil
}

func runFresh(ctx context.Context, db *pgxpool.Pool, migrations []Migration) error {
	applied, err := getAppliedMigrations(ctx, db)
	if err != nil {
		return err
	}

	for i := len(migrations) - 1; i >= 0; i-- {
		migration := migrations[i]

		if !applied[migration.Version] {
			continue
		}

		if err := rollbackMigration(ctx, db, migration); err != nil {
			return err
		}
	}

	return runUp(ctx, db, migrations)
}

func runStatus(ctx context.Context, db *pgxpool.Pool, migrations []Migration) error {
	applied, err := getAppliedMigrations(ctx, db)
	if err != nil {
		return err
	}

	for _, migration := range migrations {
		status := "pending"

		if applied[migration.Version] {
			status = "applied"
		}

		fmt.Printf("%s %-40s %s\n", migration.Version, migration.Name, status)
	}

	return nil
}

func normalizeMigrationName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "-", "_")
	return name
}

func createFileWithContent(migrationName, path, content string) error {
	filePath := filepath.Join(path, migrationName)

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		return err
	}

	return nil
}

func runCreateFile(path string, tableNames []string) error {
	version := time.Now().Format("20060102150405")

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return fmt.Errorf("erro ao criar diretório de migrations: %w", err)
	}

	for _, tableName := range tableNames {
		if strings.TrimSpace(tableName) == "" {
			return fmt.Errorf("nome da tabela é obrigatório")
		}

		// 20260707120000_create_tableName.up.sql
		// 20260707120000_create_tableName.down.sql
		tableName = normalizeMigrationName(tableName)

		upFileName := fmt.Sprintf("%s_create_%s.up.sql", version, tableName)
		downFileName := fmt.Sprintf("%s_create_%s.down.sql", version, tableName)

		upContent := fmt.Sprintf(`CREAT TABLE IF NOT EXISTS %s (
			created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			deleted_at TIMESTAMPTZ NULL
		);`, tableName)

		downContent := fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, tableName)

		if err := createFileWithContent(upFileName, path, upContent); err != nil {
			return err
		}

		if err := createFileWithContent(downFileName, path, downContent); err != nil {
			return err
		}

	}

	return nil
}

func Run(ctx context.Context, db *pgxpool.Pool, option Options) error {
	path, err := resolveMigrationDir(option.Dir)

	if err != nil {
		return err
	}

	migrations, err := collectMigrations(path)

	if err != nil {
		return err
	}

	if err := runSchemaMigrations(ctx, db); err != nil {
		return err
	}

	switch option.Command {
	case CommandUp:
		return runUp(ctx, db, migrations)
	case CommandDown:
		return runDown(ctx, db, migrations)
	case CommandFresh:
		return runFresh(ctx, db, migrations)
	case CommandStatus:
		return runStatus(ctx, db, migrations)
	case CommandCreate:
		return runCreateFile(path, option.ExtraArgs)
	default:
		return fmt.Errorf("comando de migration inválido: %s", option.Command)
	}
}
