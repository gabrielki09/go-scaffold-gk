package scaffold

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

var (
	ErrModelNameRequired = errors.New("flag m é obrigatória")
	ErrIDTypeRequired    = errors.New("informe o tipo de ID: -uuid ou -id")
	ErrOnlyOneIDType     = errors.New("somente um tipo de ID pode ser utilizado")
	ErrRootDirRequired   = errors.New("caminho principal não informado")
	ErrRootDirNotExists  = errors.New("caminho principal não existe")
)

const (
	CommandMigration      = "M"
	CommandModel          = "m"
	CommandUUIDUse        = "uuid_use"
	CommandIDUse          = "id_use"
	CommandRequests       = "requests"
	CommandResource       = "resource"
	CommandSeed           = "seed"
	CommandRepository     = "repo"
	CommandController     = "controller"
	CommandCreateRepoPath = "create_repo_path"

	CommandRoutes  = "routes"
	CommandService = "service"
)

var allowedCommands = map[string]struct{}{
	CommandModel:          {},
	CommandUUIDUse:        {},
	CommandIDUse:          {},
	CommandRequests:       {},
	CommandResource:       {},
	CommandSeed:           {},
	CommandMigration:      {},
	CommandController:     {},
	CommandRepository:     {},
	CommandCreateRepoPath: {},

	// For repo pattern
	CommandRoutes:  {},
	CommandService: {},
}

func exists(s string) error {
	_, err := os.Stat(s)
	if err == nil {
		return nil
	}

	return err
}

func (o Options) Validate() error {
	if strings.TrimSpace(o.Name) == "" {
		return ErrModelNameRequired
	}

	if o.Command == nil {
		return errors.New("command map não pode ser nil")
	}

	for command := range o.Command {
		if _, ok := allowedCommands[command]; !ok {
			return fmt.Errorf("comando inválido: %s", command)
		}
	}

	uuidUse := o.Command[CommandUUIDUse]
	idUse := o.Command[CommandIDUse]

	if !uuidUse && !idUse {
		return ErrIDTypeRequired
	}

	if uuidUse && idUse {
		return ErrOnlyOneIDType
	}

	if o.Command["repo"] {
		if !o.Command["create_repo_path"] {
			if err := exists(o.RootDir); err != nil {
				return ErrRootDirNotExists
			}
		} else {
			if err := createInformedPath(o.RootDir, 0755); err != nil {
				return err
			}
		}
	}

	return nil
}
