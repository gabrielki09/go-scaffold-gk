package model

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func buildModelFullPath(model, path string, separateByFolder bool) (string, error) {
	fileName := func(s string) string {
		return fmt.Sprintf("%s_model.go", normalizeModelFileName(s))
	}

	if !separateByFolder {
		return filepath.Join(path, fileName(model)), nil
	}

	fullPath := filepath.Join(path, strings.ToLower(model))

	exists, err := existsPath(fullPath)
	if err != nil {
		return "", err
	}

	if !exists {
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return "", err
		}

		fullPath = filepath.Join(path, strings.ToLower(model), fileName(model))
	}

	return fullPath, nil
}

func buildModelContent(model string) string {
	capitalize := func(s string) string {
		var newStr string
		isFirstChar := true

		for _, c := range strings.Split(s, "") {
			if isFirstChar {
				newStr += strings.ToUpper(c)
				isFirstChar = !isFirstChar
			} else {
				newStr += c

			}

		}

		return newStr
	}

	content := fmt.Sprintf(`package %smodel

type %sModel struct {
}
	`, model, capitalize(model))

	return content
}

func createModelFile(model, path string, separateByFolder bool) error {
	filePath, err := buildModelFullPath(model, path, separateByFolder)

	if err != nil {
		return err

	}

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	if _, err := file.WriteString(buildModelContent(normalizeModelNameContent(model))); err != nil {
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

	return createModelFile(option.ModelName, dir, option.SeparateByFolder)
}
