package scaffoldv2

func Run(options Options) error {

	if err := options.ValidateOptions(); err != nil {
		return err
	}

	dirs, err := resolveFileDir(options)
	if err != nil {
		return err
	}

	goMod, err := getModuleName()
	if err != nil {
		return err
	}

	filesConfig := FileConfig{
		Name:       options.ModelName,
		FilePaths:  dirs,
		ModuleName: goMod,
		RootDir:    options.RootDir,
	}

	return runCreateFiles(filesConfig, options)

}
