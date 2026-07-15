package scaffold

type CommandOptions struct {
	UUIDUse          bool
	IDUse            bool
	SeparateByFolder bool
	Requests         bool
	Resource         bool
	Seed             bool
	Migration        bool
	Controller       bool
	All              bool
}

func NewCommandMap(options CommandOptions) map[string]bool {
	commands := map[string]bool{
		CommandModel:            true,
		CommandUUIDUse:          options.UUIDUse,
		CommandIDUse:            options.IDUse,
		CommandSeparateByFolder: options.SeparateByFolder,
		CommandRequests:         options.Requests,
		CommandResource:         options.Resource,
		CommandSeed:             options.Seed,
		CommandMigration:        options.Migration,
		CommandController:       options.Controller,
	}

	if options.All {
		commands[CommandRequests] = true
		commands[CommandResource] = true
		commands[CommandSeed] = true
		commands[CommandMigration] = true
		commands[CommandController] = true
	}

	return commands
}
