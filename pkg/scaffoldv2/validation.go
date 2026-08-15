package scaffoldv2

import (
	"fmt"
	"os"
	"strings"

	"charm.land/log/v2"
)

func pathExists(s string) error {
	_, err := os.Stat(s)

	if err == nil {
		return nil
	}

	log.Infof("O caminho %s não foi localizado, o mesmo será criado.", s)
	if err := globalCreateInformedPath(s, 0755); err != nil {
		log.Error("Erro ao criar o caminho: ", err)
		return err
	}

	return nil
}

func (o Options) ValidateOptions() error {
	if strings.TrimSpace(o.ModelName) == "" {
		return ErrModelNameRequired
	}

	if o.CommandFlags == nil {
		return ErrNoInformedFlags
	}

	for command := range o.CommandFlags {
		if _, ok := allowedCommands[command]; !ok {
			return fmt.Errorf("comando inválido: %s", command)
		}
	}

	uuidUse := o.CommandFlags[CommandUUIDUse]
	idUse := o.CommandFlags[CommandIDUse]

	if !uuidUse && !idUse {
		return ErrIDTypeRequired
	}

	if uuidUse && idUse {
		return ErrOnlyOneIDType
	}

	if o.RootDir == "" {
		return ErrRootDirRequired
	}

	return nil
}
