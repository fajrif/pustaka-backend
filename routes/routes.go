package routes

import (
    "pustaka-backend/handlers"
    "pustaka-backend/middleware"
    "github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
    api := app.Group("/api")

    // Public routes
    auth := api.Group("/auth")
    auth.Post("/register", handlers.Register)
    auth.Post("/login", handlers.Login)

    // Protected routes
    api.Use(middleware.AuthRequired())

    // User routes
    api.Get("/me", handlers.GetMe)
    api.Put("/me", handlers.UpdateMe)

    // MerkBuku routes
    merkBuku := api.Group("/merk-buku")
    merkBuku.Get("/", handlers.GetAllMerkBuku)
    merkBuku.Get("/:id", handlers.GetMerkBuku)
    merkBuku.Post("/", handlers.CreateMerkBuku)
    merkBuku.Put("/:id", handlers.UpdateMerkBuku)
    merkBuku.Delete("/:id", handlers.DeleteMerkBuku)

}
