package env

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var AdminPassword string
var AllowOrigin string

func init() {
	// .env 파일이 없어도 괜찮음 (Render 등 클라우드에서는 OS 환경변수로 직접 주입)
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Info: .env file not found, using OS environment variables")
	}
	AdminPassword = os.Getenv("ADMIN_PASSWORD")
	AllowOrigin = os.Getenv("ALLOW_ORIGIN")
}
