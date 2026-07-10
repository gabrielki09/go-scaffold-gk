package migration

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func collectMigrations(dir string) ([]Migration, error) {
	migrationByVersion := map[string]*Migration{}

	err := filepath.WalkDir(dir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			return nil
		}

		if !strings.HasSuffix(entry.Name(), ".sql") {
			return err
		}

		version, name, direction, ok := parseMigrationFileName(entry.Name())

		if !ok {
			return fmt.Errorf("nome de migration inválido: %s", entry.Name())
		}

		migration, exists := migrationByVersion[version]
		if !exists {
			migration = &Migration{
				Version: version,
				Name:    name,
			}

			migrationByVersion[version] = migration
		}

		if migration.Name != name {
			return fmt.Errorf(
				"versão %s possui nomes divergentes: %s e %s",
				version,
				migration.Name,
				name,
			)
		}

		switch direction {
		case "up":
			migration.UpFile = path
		case "down":
			migration.DownFile = path
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	migrations := make([]Migration, 0, len(migrationByVersion))

	for _, migration := range migrationByVersion {
		if migration.UpFile == "" {
			return nil, fmt.Errorf("migration %s_%s não possui arquivo .up.sql", migration.Version, migration.Name)
		}

		if migration.DownFile == "" {
			return nil, fmt.Errorf("migration %s_%s não possui arquivo .down.sql", migration.Version, migration.Name)
		}

		migrations = append(migrations, *migration)
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

func resolveMigrationDir(dir string) (string, error) {
	if strings.TrimSpace(dir) != "" {
		return filepath.Abs(dir)
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Join(wd, "database", "migration"), nil
}
