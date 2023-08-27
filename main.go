package main

import (
	"github.com/dotenv-org/godotenvvault"
	"github.com/saigees/michelangelo/internal"
	"github.com/saigees/michelangelo/pkg/michelangelo"
)

func main() {
	err := godotenvvault.Load()
	if err != nil {
		internal.Logger.Error("Failed to load .env")
	}

	michelangelo.Main()
}
