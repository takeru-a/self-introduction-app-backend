package configs

import (
    "log"
    "os"
    "github.com/joho/godotenv"
)

// 環境変数をロードする
func EnvMongoURI() string {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    return os.Getenv("MONGOURI")
}