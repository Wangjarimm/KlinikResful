package main

import (
    "context"
    "fmt"
    "log"
    "os/signal"
    "sistem-informasi-klinik/database"
    "sistem-informasi-klinik/router"
    "syscall"
    "time"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/joho/godotenv"
)

// app adalah variabel global untuk menyimpan instance Fiber app
var app *fiber.App

func init() {
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }
}

func main() {
    // Initialize Fiber app
    app = fiber.New()

    // Connect to database
    database.ConnectDB()

    // Setup middleware
    app.Use(cors.New(cors.Config{
        AllowHeaders:     "*",
        AllowCredentials: true,
        AllowMethods:     "GET,POST,PUT,DELETE",
    }))

    // Setup routes
    router.SetupRoutes(app)

    // Create a context that will be canceled when a signal is caught
    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer stop()

    // Start the server in a goroutine
    go func() {
        if err := app.Listen(":4000"); err != nil {
            log.Fatalf("Error starting server: %v\n", err)
        }
    }()

    // Log server start time
    fmt.Println("Server started at", time.Now().Format("2006-01-02 15:04:05"))

    // Wait for a signal to shut down
    <-ctx.Done()

    // Stop receiving any more signals
    stop()

    // Gracefully shut down the server
    fmt.Println("Shutting down server...")
    if err := app.Shutdown(); err != nil {
        log.Printf("Error shutting down server: %v\n", err)
    } else {
        fmt.Println("Server shut down gracefully.")
    }
}
