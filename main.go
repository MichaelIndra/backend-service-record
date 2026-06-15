package main

import (
	"fmt"
	"os"
	"service-record/config"
	"service-record/controllers"
	"service-record/docs" // Ini penting untuk Swagger

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Service Kendaraan API
// @version 1.0
// @description API untuk mencatat servis kendaraan dan manajemen foto per transaksi.
// @BasePath /
func main() {
	// 1. Inisialisasi Database
	config.ConnectDatabase()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		// Fungsi ini akan memeriksa setiap request yang masuk
		AllowOriginFunc: func(origin string) bool {
			// Opsi A: Izinkan SEMUA origin secara dinamis (Sama seperti "*", tapi mendukung AllowCredentials)
			return true

			// Opsi B: Atau beri validasi tertentu, misal hanya domain yang mengandung kata 'bengkel'
			// return strings.Contains(origin, "bengkel")
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	// 2. Akses Folder Foto (Static)
	// Agar folder uploads bisa diakses via browser (contoh: localhost:8080/uploads/trx_1/nota.jpg)
	r.Static("/uploads", "./uploads")

	docs.SwaggerInfo.Host = ""
	// 3. Swagger Route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// 4. Grouping Routes
	api := r.Group("/api")
	{
		// Transaction Routes
		api.GET("/transaction", controllers.GetTransaction)
		api.GET("/transaction/:id", controllers.GetTransactionByID)
		api.POST("/transaction", controllers.CreateTransaction)
		api.DELETE("/transaction/:id", controllers.DeleteTransaction)

		// Master Kendaraan
		api.GET("/kendaraan", controllers.GetKendaraan)
		api.GET("/kendaraan/:no_polisi", controllers.GetKendaraanByID)
		api.GET("/kendaraan/getKendaraanByStatus/:status", controllers.GetKendaraanByStatus)
		api.POST("/kendaraan", controllers.CreateKendaraan)
		api.PUT("/kendaraan/:no_polisi", controllers.UpdateKendaraan)

		// Master Bengkel
		api.GET("/bengkel", controllers.GetBengkel)
		api.GET("/bengkel/:id", controllers.GetBengkelByID)
		api.GET("/bengkel/getBengkelByStatus/:status", controllers.GetBengkelByStatus)
		api.POST("/bengkel", controllers.CreateBengkel)
		api.PUT("/bengkel/:id", controllers.UpdateBengkel)

	}

	fmt.Println("Server running on port:", port)
	r.Run(":" + port)
}
