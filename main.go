package main

import (
    "bruno_authentication/bootstrap"

    "github.com/joho/godotenv"
)

func main() {
    _ = godotenv.Load()
    _ = bootstrap.RootApp.Execute()
}
