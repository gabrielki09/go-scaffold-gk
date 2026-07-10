package migration

import "strings"

func parseMigrationFileName(fileName string) (version, name, direction string, ok bool) {
	if !strings.HasSuffix(fileName, ".sql") {
		return "", "", "", false
	}

	trimmed := strings.TrimSuffix(fileName, ".sql")

	parts := strings.Split(trimmed, ".")
	if len(parts) != 2 {
		return "", "", "", false
	}

	baseName := parts[0]
	direction = parts[1]

	if direction != "up" && direction != "down" {
		return "", "", "", false
	}

	nameParts := strings.SplitN(baseName, "_", 2)
	if len(nameParts) != 2 {
		return "", "", "", false
	}

	version = nameParts[0]
	name = nameParts[1]

	return version, name, direction, true

}
