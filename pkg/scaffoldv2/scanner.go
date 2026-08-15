package scaffoldv2

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func resolveFileDir(option Options) (map[string]string, error) {
	var rootPath = option.RootDir
	allPaths := make(map[string]string)

	for key, defaultRootModuleDir := range moduleRootDirs {

		var fullPath string

		if isSandBox {
			if (key == "m" && option.CommandFlags[CommandModel]) || (key == "M" && option.CommandFlags[CommandMigration]) {
				fullPath = filepath.Join("sandbox", defaultRootModuleDir)

			} else {
				fullPath = filepath.Join("sandbox", rootPath, defaultRootModuleDir)
			}
		} else {
			if (key == "m" && option.CommandFlags[CommandModel]) || (key == "M" && option.CommandFlags[CommandMigration]) {
				fullPath = defaultRootModuleDir
			} else {
				fullPath = filepath.Join(rootPath, defaultRootModuleDir)
			}
		}

		validatedPath, err := validatePath(fullPath)

		if err != nil {
			return nil, fmt.Errorf("erro ao validar o diretório de %s: %w", key, err)
		}

		allPaths[key] = validatedPath
	}

	return allPaths, nil
}

func getModuleName() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	fileByte, err := os.ReadFile(filepath.Join(wd, "go.mod"))
	if err != nil {
		return "", err
	}

	firstLine := strings.TrimSpace(string(bytes.SplitN(fileByte, []byte("\n"), 2)[0]))
	if firstLine == "" {
		return "your_module", nil
	}

	parts := strings.Fields(firstLine)

	if len(parts) < 2 || parts[0] != "module" {
		return "your_module", nil
	}

	return parts[1], nil
}
