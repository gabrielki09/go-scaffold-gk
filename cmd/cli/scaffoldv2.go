package main

import (
	"flag"

	"charm.land/log/v2"

	"github.com/gabrielki09/go-scaffold-gk/pkg/scaffoldv2"
)

func init() {
	log.SetReportCaller(true)
}

func main() {
	var (
		path      = flag.String("path", "", "Comando para informar o caminho principal das operações")
		modelFlag = flag.String("m", "", "Comando para criação de arquivo padrão da model")
		migration = flag.Bool("M", false, "Comando para criação da migration")
		uuidUse   = flag.Bool("uuid", false, "Comando para criação do model com uuid")
		idUse     = flag.Bool("id", false, "Comando para criação do model com id (int)")

		debug = flag.Bool("debug", false, "Executa em debug")
	)

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	flag.Parse()

	flags := make(map[string]bool)

	flags["m"] = true
	flags["uuid_use"] = *uuidUse
	flags["id_use"] = *idUse
	flags["M"] = *migration

	options := scaffoldv2.Options{
		ModelName:    *modelFlag,
		CommandFlags: flags,
		RootDir:      *path,
	}

	if err := scaffoldv2.Run(options); err != nil {
		log.Fatal(err)
	}

	log.Info("Finalizado.")
}
