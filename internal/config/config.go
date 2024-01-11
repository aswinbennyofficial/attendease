package config

import(
	"github.com/joho/godotenv"
	"log"
)


func LoadEnv(){
	err:=godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading environment variables file")
		return
	}
}