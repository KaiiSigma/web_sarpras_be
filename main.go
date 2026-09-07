package main

import (
	"fmt"
	"log"
	"os"

	"web_sarpras/config"
	"web_sarpras/controllers"
	"web_sarpras/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Println("Info: File .env tidak ditemukan, menggunakan environment OS bawaan")
	}

	config.ConnectDatabase()
	fmt.Println("Koneksi & tes ke MariaDB berhasil!")

	r := gin.Default()

	r.SetTrustedProxies(nil)

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	api := r.Group("/api")
	{
		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)

		adminRoutes := api.Group("/admin")
		adminRoutes.Use(middleware.AuthMiddleware("admin"))
		{
			adminRoutes.POST("/users", controllers.CreateUserByAdmin)

			adminRoutes.POST("/aset", controllers.CreateAset)
			adminRoutes.PUT("/aset/:id", controllers.UpdateAset)
			adminRoutes.DELETE("/aset/:id", controllers.DeleteAset)
		}

		userRoutes := api.Group("/user")
		userRoutes.Use(middleware.AuthMiddleware("user", "admin", "teknisi"))
		{
			userRoutes.GET("/aset", controllers.GetAllAset)
			userRoutes.GET("/aset/:id", controllers.GetAsetByID)
		}

		teknisiRoutes := api.Group("/teknisi")
		teknisiRoutes.Use(middleware.AuthMiddleware("teknisi", "admin"))
		{
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	err := r.Run("0.0.0.0:" + port)
	if err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}