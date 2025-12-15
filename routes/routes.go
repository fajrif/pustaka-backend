package routes

import (
    "pustaka-backend/handlers"
    "pustaka-backend/middleware"
    "github.com/gofiber/fiber/v2"

    fiberSwagger "github.com/swaggo/fiber-swagger"
)

func Setup(app *fiber.App) {
    // Swagger documentation route
    app.Get("/swagger/*", fiberSwagger.WrapHandler)

    api := app.Group("/api")

    // Public routes
    auth := api.Group("/auth")
    auth.Post("/login", handlers.Login)

    // Protected routes
    api.Use(middleware.AuthRequired())

    // User routes
    api.Get("/me", handlers.GetMe)
    api.Put("/me", handlers.UpdateMe)

    // Upload routes
    api.Post("/upload/:resource/:field/:id", handlers.UploadResourceField)
    api.Delete("/upload/:resource/:field/:id", handlers.DeleteResourceField)

    // Cities routes
    cities := api.Group("/cities")
    cities.Get("/", handlers.GetAllCities)
    cities.Get("/:id", handlers.GetCity)
    cities.Post("/", handlers.CreateCity)
    cities.Put("/:id", handlers.UpdateCity)
    cities.Delete("/:id", handlers.DeleteCity)

    // Expeditions routes
    expeditions := api.Group("/expeditions")
    expeditions.Get("/", handlers.GetAllExpeditions)
    expeditions.Get("/:id", handlers.GetExpedition)
    expeditions.Post("/", handlers.CreateExpedition)
    expeditions.Put("/:id", handlers.UpdateExpedition)
    expeditions.Delete("/:id", handlers.DeleteExpedition)

    // MerkBuku routes
    merkBuku := api.Group("/merk-buku")
    merkBuku.Get("/", handlers.GetAllMerkBuku)
    merkBuku.Get("/:id", handlers.GetMerkBuku)
    merkBuku.Post("/", handlers.CreateMerkBuku)
    merkBuku.Put("/:id", handlers.UpdateMerkBuku)
    merkBuku.Delete("/:id", handlers.DeleteMerkBuku)

    // JenisBuku routes
    jenisBuku := api.Group("/jenis-buku")
    jenisBuku.Get("/", handlers.GetAllJenisBuku)
    jenisBuku.Get("/:id", handlers.GetJenisBuku)
    jenisBuku.Post("/", handlers.CreateJenisBuku)
    jenisBuku.Put("/:id", handlers.UpdateJenisBuku)
    jenisBuku.Delete("/:id", handlers.DeleteJenisBuku)

    // JenjangStudi routes
    jenjangStudi := api.Group("/jenjang-studi")
    jenjangStudi.Get("/", handlers.GetAllJenjangStudi)
    jenjangStudi.Get("/:id", handlers.GetJenjangStudi)
    jenjangStudi.Post("/", handlers.CreateJenjangStudi)
    jenjangStudi.Put("/:id", handlers.UpdateJenjangStudi)
    jenjangStudi.Delete("/:id", handlers.DeleteJenjangStudi)

    // BidangStudi routes
    bidangStudi := api.Group("/bidang-studi")
    bidangStudi.Get("/", handlers.GetAllBidangStudi)
    bidangStudi.Get("/:id", handlers.GetBidangStudi)
    bidangStudi.Post("/", handlers.CreateBidangStudi)
    bidangStudi.Put("/:id", handlers.UpdateBidangStudi)
    bidangStudi.Delete("/:id", handlers.DeleteBidangStudi)

    // Kelas routes
    kelas := api.Group("/kelas")
    kelas.Get("/", handlers.GetAllKelas)
    kelas.Get("/:id", handlers.GetKelas)
    kelas.Post("/", handlers.CreateKelas)
    kelas.Put("/:id", handlers.UpdateKelas)
    kelas.Delete("/:id", handlers.DeleteKelas)

    // Publishers routes
    publishers := api.Group("/publishers")
    publishers.Get("/", handlers.GetAllPublishers)
    publishers.Get("/:id", handlers.GetPublisher)
    publishers.Post("/", handlers.CreatePublisher)
    publishers.Put("/:id", handlers.UpdatePublisher)
    publishers.Delete("/:id", handlers.DeletePublisher)

    // Books routes
    books := api.Group("/books")
    books.Get("/", handlers.GetAllBooks)
    books.Get("/:id", handlers.GetBook)
    books.Post("/", handlers.CreateBook)
    books.Put("/:id", handlers.UpdateBook)
    books.Delete("/:id", handlers.DeleteBook)

    // SalesAssociates routes
    salesAssociates := api.Group("/sales-associates")
    salesAssociates.Get("/", handlers.GetAllSalesAssociates)
    salesAssociates.Get("/:id", handlers.GetSalesAssociate)
    salesAssociates.Post("/", handlers.CreateSalesAssociate)
    salesAssociates.Put("/:id", handlers.UpdateSalesAssociate)
    salesAssociates.Delete("/:id", handlers.DeleteSalesAssociate)

    // SalesTransactions routes
    salesTransactions := api.Group("/sales-transactions")
    salesTransactions.Get("/", handlers.GetAllSalesTransactions)
    salesTransactions.Get("/:id", handlers.GetSalesTransaction)
    salesTransactions.Post("/", handlers.CreateSalesTransaction)
    salesTransactions.Put("/:id", handlers.UpdateSalesTransaction)
    salesTransactions.Delete("/:id", handlers.DeleteSalesTransaction)
    salesTransactions.Get("/:transaction_id/installments", handlers.GetTransactionInstallments)
    salesTransactions.Post("/:transaction_id/installments", handlers.AddInstallment)
		salesTransactions.Delete("/:transaction_id/installments/:id", handlers.DeleteInstallment)

    // Billers routes
    billers := api.Group("/billers")
    billers.Get("/", handlers.GetAllBillers)
    billers.Get("/:id", handlers.GetBiller)
    billers.Post("/", handlers.CreateBiller)
    billers.Put("/:id", handlers.UpdateBiller)
    billers.Delete("/:id", handlers.DeleteBiller)

    // Admin Only routes
    api.Use(middleware.AdminOnly())

    // Users routes
    users := api.Group("/users")
    users.Get("/", handlers.GetAllUsers)
    users.Get("/:id", handlers.GetUser)
    users.Post("/", handlers.CreateUser)
    users.Put("/:id", handlers.UpdateUser)
    users.Delete("/:id", handlers.DeleteUser)
}
