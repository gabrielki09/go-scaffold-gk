package scaffoldv2

type Options struct {
	ModelName    string
	CommandFlags map[string]bool
	RootDir      string
}

type FileConfig struct {
	ModuleName string
	Name       string
	RootDir    string
	FilePaths  map[string]string
}
