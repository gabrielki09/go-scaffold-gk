package main

import (
	"flag"
	"log"

	"github.com/gabrielki09/go-migrate-gk/pkg/model"
)

func main() {
	var (
		modelFlag        = flag.String("model", "", "Comando para criação de arquivo padrão da model")
		separateByFolder = flag.Bool("s", false, "Comando para separação de pastas por model")
		makeAll          = flag.Bool("a", false, "Comando para separação de pastas por model")
	)

	flag.Parse()

	if *modelFlag != "" {
		log.Println(*modelFlag)
		if err := model.Run(model.Options{
			ModelName:        *modelFlag,
			SeparateByFolder: *separateByFolder,
		}); err != nil {
			log.Fatal(err)
		}
	}
}
