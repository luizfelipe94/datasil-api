package main

import (
	"fmt"
	"log"

	"github.com/luizfelipe94/datasil/cmd/api"
	"github.com/luizfelipe94/datasil/configs"
	"github.com/luizfelipe94/datasil/db"
)

const PORT = "5000"

func main() {
	db := db.NewDB(&configs.Envs)
	server := api.NewAPIServer(fmt.Sprintf(":%s", PORT), db)
	log.Println("Server running on port", PORT)
	if err := server.Run(); err != nil {
		panic(err)
	}
}
