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

func GetSIGNING_KEY() string {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
    return os.Getenv("SIGNING_KEY")
}

func GetRedisEndpoint() string {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
    return os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")
}

func GetSESSION_KEY() string {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
    return os.Getenv("SESSION_KEY")
}
 