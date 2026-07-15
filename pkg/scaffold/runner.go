package scaffold

func Run(option Options) error {
	if err := option.Validate(); err != nil {
		return err
	}

	dirs, err := resolveFileDir(option.Command)
	if err != nil {
		return err
	}

	fileConfig := File{
		Name:             option.Name,
		FilePaths:        dirs,
		SeparateByFolder: option.SeparateByFolder,
	}

	return createFiles(fileConfig, option)

}
