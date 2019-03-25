package main

import (
	"github.com/gin-gonic/gin"
	"os"
	"log"
	"github.com/joho/godotenv"
	"./api"
	"./database"
	"./lib/middlewares"
)

func main() {
  err := godotenv.Load()
	if err != nil {
		log.Fatal("Error al cargar el archivo .env")
	}

	db, _ := database.Initialize() // iniciando conexion con la base de datos
	
	port := os.Getenv("PORT")
	app := gin.Default()					// creando la app gin
	app.Use(database.Inject(db))
	app.Use(middlewares.JWTMiddleware()) // validando rutas
	api.ApplyRoutes(app)	// aplicando api al router
	app.Run(":" + port)						// escuchando el el puerto
}
