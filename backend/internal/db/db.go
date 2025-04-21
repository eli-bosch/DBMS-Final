package db

import (
	"fmt"

	"github.com/eli-bosch/DBMS-final/internal/models"

	"github.com/eli-bosch/DBMS-final/config"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func init() {

	db := config.Connect()

	db.AutoMigrate(models.Assignment{})
	db.AutoMigrate(models.Building{})
	db.AutoMigrate(models.Room{})
	db.AutoMigrate(models.Student{})

	fmt.Println("Database is connected...")
}
