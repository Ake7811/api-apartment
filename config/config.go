package config

import (
	"apartment/api/security"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	PORT     = 0
	DBURL    = ""
	DBDRIVER = ""
)

func Load() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env: %v", err)
	}
	PORT, err = strconv.Atoi(os.Getenv("API_PORT"))
	if err != nil {
		PORT = 7000
	}
	key := hex.EncodeToString([]byte(os.Getenv("KEY")))
	decrypted := security.Decrypt(os.Getenv("DB_PASS"), key)

	DBDRIVER = os.Getenv("DB_DRIVER")
	DBURL = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True", os.Getenv("DB_USER"), decrypted, os.Getenv("DB_IP"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
}
