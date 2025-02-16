package main

import (
	"fmt"
	"github.com/btnbrd/avitoshop/internal/application"
	"github.com/btnbrd/avitoshop/internal/storage"
	"log"
)

func main() {
	db, err := storage.NewDBConnection()
	fmt.Println("err=", err)
	if err != nil {
		log.Fatalf(err.Error())
		return
	}
	defer db.Close()

	apiServer := application.NewServer(db)
	if err := apiServer.Run(); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
		//panic(err)
	}
}
