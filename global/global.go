// global - global variables
package global

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	ADMIN_PASSWORD string
	PORT           string
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
		e.Write([]byte("PORT=\nADMIN_PASSWORD=\n"))
		log.Println("Please initialize variables in .env")
		os.Exit(1)
	}
	ADMIN_PASSWORD = os.Getenv("ADMIN_PASSWORD")
	PORT = os.Getenv("PORT")
}
