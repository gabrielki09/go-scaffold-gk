package scaffold

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrModelNameRequired = errors.New("flag m é obrigatória")
	ErrIDTypeRequired    = errors.New("informe o tipo de ID: -uuid ou -id")
	ErrOnlyOneIDType     = errors.New("somente um tipo de ID pode ser utilizado")
)

const (
	CommandModel            = "m"
	CommandUUIDUse          = "uuid_use"
	CommandIDUse            = "id_use"
	CommandSeparateByFolder = "separate_by_folder"
	CommandRequests         = "requests"
	CommandResource         = "resource"
	CommandSeed             = "seed"
	CommandMigration        = "migration"
	CommandController       = "controller"
)

var allowedCommands = map[string]struct{}{
	CommandModel:            {},
	CommandUUIDUse:          {},
	CommandIDUse:            {},
	CommandSeparateByFolder: {},
	CommandRequests:         {},
	CommandResource:         {},
	CommandSeed:             {},
	CommandMigration:        {},
	CommandController:       {},
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

	if uuidUse && idUse {
		return ErrOnlyOneIDType
	}

	if !uuidUse && !idUse {
		return ErrIDTypeRequired
	}

	return nil
}
