package model

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func runCreateModelFile(model, path string, separateByFolder bool) error {
	var filePath string

	normalizedModelFileName := normalizeModelFileName(model)
	if !separateByFolder {
		filePath = filepath.Join(path, fmt.Sprintf("%s_model.go", normalizedModelFileName))

	} else {

		fullPath := filepath.Join(path, normalizedModelFileName)
		exists, err := existsPath(fullPath)
		if err != nil {
			return err
		}

		if !exists {
			if err := os.MkdirAll(fullPath, 0755); err != nil {
				return err
			}
		}

		filePath = fmt.Sprintf("%s/%s_model.go", fullPath, normalizedModelFileName)
	}

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	if _, err := file.WriteString(fmt.Sprintf(`package %smodel`, normalizeModelNameContent(model))); err != nil {
		return err
	}

	return nil
}

// Input: Financial Account
// Output: financial_account
// Use para o nome que será dado ao arquivo
func normalizeModelFileName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "-", "_")
	return name
}

// Input: Financial Account
// Output: financialaccount
// Use para o conteúdo que será escrito no arquivo
func normalizeModelNameContent(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "")
	name = strings.ReplaceAll(name, "-", "")
	return name
}

func validateModelName(modelName string) error {
	if modelName == "" {
		return fmt.Errorf("o nome do model é obrigatório.")
	}

	if _, err := strconv.Atoi(modelName); err == nil {
		return fmt.Errorf("o nome do model deve ser uma string.")
	}

	return nil
}

func Run(option Options) error {
	if err := validateModelName(option.ModelName); err != nil {
		return err
	}

	dir, err := resolveModelDir()

	if err != nil {
		return err
	}

	return runCreateModelFile(option.ModelName, dir, option.SeparateByFolder)
}
