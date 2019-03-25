package database

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"./models"
)

// Iniciando la base de datos
func Initialize() (*gorm.DB, error) {
	dbConfig := os.Getenv("DB_CONFIG")
	db, err := gorm.Open("mysql", dbConfig)
	db.LogMode(true)
	if err != nil {
		fmt.Println("Ocurrio un error en la conexión con la base de datos =>", err)
		panic(err)
	}
	fmt.Println("Connected to database")
	models.Migrate(db)
	return db, err
}

