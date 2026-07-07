package migration

type Command string

const (
	CommandUp     Command = "up"
	CommandDown   Command = "down"
	CommandFresh  Command = "fresh"
	CommandStatus Command = "status"
	CommandCreate Command = "create"
)

type Options struct {
	Dir       string
	Command   Command
	ExtraArgs []string
}

type Migration struct {
	Version  string
	Name     string
	UpFile   string
	DownFile string
}
