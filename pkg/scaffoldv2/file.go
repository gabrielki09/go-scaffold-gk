package scaffoldv2

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"charm.land/log/v2"
)

func buildMigrationContent(fileName, usageId string) (
	upFileName,
	downFileName,
	upContent,
	downContent string,
	err error,
) {
	version := time.Now().Format("20060102150405")

	upFileName = fmt.Sprintf("%s_create_%s.up.sql", version, fileName)
	downFileName = fmt.Sprintf("%s_create_%s.down.sql", version, fileName)

	upContent = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
	%s,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	deleted_at TIMESTAMPTZ NULL
);`, fileName, usageId)

	downContent = fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, fileName)

	return upFileName, downFileName, upContent, downContent, nil
}

func createMigrationFile(fileName, filePath string, usageId bool) error {
	var migrationId string

	if usageId {
		migrationId = "id BIGSERIAL PRIMARY KEY"
	} else {
		migrationId = "id UUID PRIMARY KEY DEFAULT gen_random_uuid()"
	}

	migrationUpFileName, migrationDownFileName, migrationUpContent, migrationDownContent, err := buildMigrationContent(fileName, migrationId)
	if err != nil {
		return err
	}

	if err := createRawFile(migrationUpFileName, migrationUpContent, filePath); err != nil {
		return fmt.Errorf("erro ao criar o arquivo .up da migration: %w", err)
	}

	if err := createRawFile(migrationDownFileName, migrationDownContent, filePath); err != nil {
		return fmt.Errorf("erro ao criar o arquivo .down da migration: %w", err)
	}

	return nil
}

func createRawFile(fileName, fileContent, filePath string) error {
	fullPath := filepath.Join(filePath, fileName)

	osFile, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			log.Info("Arquivo já existente.")
			return nil
		}
		return err
	}

	defer osFile.Close()

	if _, err := osFile.WriteString(fileContent); err != nil {
		return err
	}

	return nil
}

func createFileTextTemplate(fileName, fileModule, filePath, getModuleName string, isID bool) error {
	fileContent, err := createContent(globalTrimSpace(fileName), fileModule, getModuleName, isID)
	if err != nil {
		return err
	}

	goFileName := fmt.Sprintf("%s_%s.go", fileName, moduleEquivalence(fileModule))

	fullPath := filepath.Join(filePath, goFileName)

	osFile, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			log.Info("Arquivo já existente.")
			return nil
		}

		return err
	}

	defer osFile.Close()

	if _, err := osFile.WriteString(fileContent); err != nil {
		return err
	}

	return nil
}

func runCreateFiles(fileConfig FileConfig, option Options) error {
	log.Info("Criando arquivos ...")

	for key, path := range fileConfig.FilePaths {
		if key == "M" {
			if !option.CommandFlags[CommandModel] {
				continue
			}

			if err := createMigrationFile(fileConfig.Name, path, option.CommandFlags[CommandIDUse]); err != nil {
				return err
			}

			continue
		}

		if err := createFileTextTemplate(
			fileConfig.Name,
			key,
			path,
			fileConfig.ModuleName,
			option.CommandFlags[CommandIDUse],
		); err != nil {
			return err
		}
	}

	return nil
}
