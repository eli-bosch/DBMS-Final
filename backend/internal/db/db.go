package db

import (
	"fmt"
	"os"

	"github.com/eli-bosch/DBMS-final/internal/models"

	"github.com/eli-bosch/DBMS-final/config"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func init() {

	db := config.Connect()

	// In dev/Turing only: drop every table so AutoMigrate can rebuild them
	if os.Getenv("FORCE_SCHEMA") == "true" {
		fmt.Println("Dropping all tables for fresh schema")
		// You can list them all in one call:
		db.DropTable(
		  &models.Assignment{},
		  &models.Room{},
		  &models.Student{},
		  &models.Building{},
		)
	  } 

	db.AutoMigrate(
		&models.Assignment{},
		&models.Building{},
		&models.Room{},
		&models.Student{},
	)

	fmt.Println("Database is connected...")
}
