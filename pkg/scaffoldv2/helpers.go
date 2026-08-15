package scaffoldv2

import (
	"os"
	"path/filepath"
	"strings"

	"charm.land/log/v2"
)

func globalCreateInformedPath(path string, perm os.FileMode) error {
	log.Infof("Caminho %s está sendo criado", path)

	if err := os.MkdirAll(path, perm); err != nil {
		log.Error("Erro ao criar o caminho: ", err)
		return err
	}

	log.Infof("Caminho %s criado com sucesso!", path)
	return nil
}

func validatePath(path string) (string, error) {
	fullPath, err := filepath.Abs(path)

	if err != nil {
		return "", err
	}

	if err := pathExists(fullPath); err != nil {
		if err := globalCreateInformedPath(fullPath, 0755); err != nil {
			return "", err
		}
	}

	return fullPath, nil
}

func globalTrimSpace(value string) string {
	return strings.TrimSpace(value)
}

func globalToLower(value string) string {
	return strings.ToLower(value)
}

func invertBarPath(value string) string {
	return strings.ReplaceAll(value, "\\", "/")
}

// Input: Financial Account
// Output: financial_account
func globalNormalizeWithUnderline(name string) string {
	name = globalTrimSpace(name)
	name = globalToLower(name)
	name = strings.ReplaceAll(name, "-", "_")
	name = strings.ReplaceAll(name, " ", "_")

	return name
}

// Input: Financial Account
// Output: financialaccount
func globalNormalizeNoWithUnderline(name string) string {
	name = globalTrimSpace(name)
	name = globalToLower(name)
	name = strings.ReplaceAll(name, "-", "")
	name = strings.ReplaceAll(name, " ", "")
	name = strings.ReplaceAll(name, "_", "")

	return name
}

func globalToPascalCase(value string) string {
	value = globalTrimSpace(value)
	value = strings.ToLower(value)
	value = strings.ReplaceAll(value, "-", "_")
	value = strings.ReplaceAll(value, " ", "_")

	parts := strings.Split(value, "_")

	var builder strings.Builder

	for _, part := range parts {
		part = globalTrimSpace(part)

		if part == "" {
			continue
		}

		builder.WriteString(strings.ToUpper(part[:1]))
		builder.WriteString(part[1:])
	}

	return builder.String()
}

func globalToKebabCase(value string) string {
	value = globalTrimSpace(value)
	value = globalToLower(value)
	value = strings.ReplaceAll(value, " ", "-")
	value = strings.ReplaceAll(value, "_", "-")

	return value
}

func moduleEquivalence(m string) string {
	// m - model
	// r - rotues
	// c - controller
	// s - service
	// R - repository
	switch m {
	case "m":
		return "model"
	case "r":
		return "routes"
	case "c":
		return "controller"
	case "s":
		return "service"
	case "R":
		return "repository"
	default:
		return ""
	}
}

func buildIDParamName(fileName string) string {
	pascalName := string(buildStructName(fileName))

	if pascalName == "" {
		return "id"
	}

	return strings.ToLower(pascalName[:1]) + pascalName[1:] + "ID"
}
