// global - global variables
package global

import (
	"log"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

var (
	HOST           string
	PORT           string
	ADMIN_EMAIL    string
	ADMIN_PASSWORD string
	AUTH_EMAIL     string
	AUTH_PASSWORD  string
	AUTH_HOST      string
	AUTH_PORT      string
	Auth           smtp.Auth
	ActiveTab      map[int]int64
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Creating .env file")
		e, err := os.Create(".env")
		if err != nil {
			log.Fatal(err)
		}
		defer e.Close()
		e.Write([]byte("HOST=\nPORT=\nADMIN_PASSWORD=\nADMIN_EMAIL=\nAUTH_EMAIL=\nAUTH_PASSWORD=\nAUTH_HOST=\nAUTH_PORT="))
		log.Println("Please initialize variables in .env")
		os.Exit(1)
	}
	ADMIN_EMAIL = os.Getenv("ADMIN_EMAIL")
	ADMIN_PASSWORD = os.Getenv("ADMIN_PASSWORD")
	AUTH_EMAIL = os.Getenv("AUTH_EMAIL")
	AUTH_PASSWORD = os.Getenv("AUTH_PASSWORD")
	AUTH_HOST = os.Getenv("AUTH_HOST")
	AUTH_PORT = os.Getenv("AUTH_PORT")
	HOST = os.Getenv("HOST")
	Auth = smtp.PlainAuth("", AUTH_EMAIL, AUTH_PASSWORD, AUTH_HOST)
	ActiveTab = make(map[int]int64)
}
