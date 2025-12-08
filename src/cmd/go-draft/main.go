package main

import (
	"log"

	"github.com/nobuenhombre/go-draft/src/cmd/go-draft/di"
)

func main() {
	// Инициализация через Wire (все зависимости создаются здесь)
	app, cleanup, err := di.InitializeApp()
	if err != nil {
		log.Fatalf("Ошибка инициализации: %v", err)
	}
	defer cleanup()

	err = app.Run()
	if err != nil {
		log.Fatalf("Ошибка приложения: %v", err)
	}
}
