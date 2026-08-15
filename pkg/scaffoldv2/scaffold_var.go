package scaffoldv2

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

var isSandBox = false

func init() {
	_ = godotenv.Load()

	if globalToLower(os.Getenv("APP_MODE")) == "debug" {
		isSandBox = !isSandBox

		if err := pathExists("sandbox"); err != nil {
			if err := globalCreateInformedPath("sandbox", 0755); err != nil {
				return
			}
		}
	}
}

const (
	DefaultRootMigrationDir  = "migration"
	DefaultRootModelDir      = "models"
	DefaultRootRoutesDir     = "routes"
	DefaultRootControllerDir = "controller"
	DefaultRootServiceDir    = "service"
	DefaultRootRepositoryDir = "repository"
)

const (
	CommandModel      = "m"
	CommandMigration  = "M"
	CommandUUIDUse    = "uuid_use"
	CommandIDUse      = "id_use"
	CommandPath       = "path"
	CommandCreatePath = "create-path"
)

var allowedCommands = map[string]struct{}{
	CommandModel:      {},
	CommandMigration:  {},
	CommandUUIDUse:    {},
	CommandIDUse:      {},
	CommandPath:       {},
	CommandCreatePath: {},
}

// R = Repository
// r = Routes
var moduleRootDirs = map[string]string{
	CommandModel:     DefaultRootModelDir,
	CommandMigration: DefaultRootMigrationDir,
	"r":              DefaultRootRoutesDir,
	"c":              DefaultRootControllerDir,
	"s":              DefaultRootServiceDir,
	"R":              DefaultRootRepositoryDir,
}

var (
	ErrModelNameRequired = errors.New("nome do model é obrigatória.")
	ErrIDTypeRequired    = errors.New("informe o tipo de ID: -uuid ou -id.")
	ErrOnlyOneIDType     = errors.New("somente um tipo de ID pode ser utilizado.")
	ErrRootDirRequired   = errors.New("caminho principal não informado.")
	ErrRootDirNotExists  = errors.New("caminho principal não existe.")
	ErrNoInformedFlags   = errors.New("nem uma flag informada.")
)
