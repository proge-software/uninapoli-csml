package envext

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

//LoadDotenvs loads the environment variables from .env files
func LoadDotenvs(appEnvKey string) {
	env := GetEnvOrDefault(appEnvKey, "development")
	godotenv.Load(".env." + env + ".local")
	if "test" != env {
		godotenv.Load(".env.local")
	}

	envfile := ".env." + env
	if err := godotenv.Load(envfile); err != nil {
		cwd, _ := os.Getwd()
		log.Printf("error loading file %v (%s): %v", envfile, cwd, err)
	}

	if err := godotenv.Load(); err != nil { // The Original .env
		cwd, _ := os.Getwd()
		log.Printf("error loading .env (%v) file: %v", cwd, err)
	}
}
