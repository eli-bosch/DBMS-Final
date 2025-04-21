package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/eli-bosch/DBMS-final/config"
	_ "github.com/eli-bosch/DBMS-final/internal/db"
	"github.com/eli-bosch/DBMS-final/internal/models"
	"github.com/eli-bosch/DBMS-final/internal/routes"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	routes.DBMSRoutes(r)

	http.Handle("/", r)
	fmt.Println("Server is running...")
	log.Fatal(http.ListenAndServe("localhost:9010", r))

}

func populateDatabase() {
	db := config.GetDB()
	rand.Seed(time.Now().UnixNano())

	// Dorm buildings at University of Arkansas
	buildingNames := []string{
		"Maple Hill South",
		"Maple Hill East",
		"Maple Hill West",
		"Founders Hall",
		"Hotz Honors Hall",
	}

	// Create buildings
	var buildings []models.Building
	for i, name := range buildingNames {
		building := models.Building{
			BuildingID:   uint(i + 1),
			BuildingName: name,
			HasAC:        rand.Intn(2) == 1,
			HasDining:    rand.Intn(2) == 1,
		}
		buildings = append(buildings, building)
		db.Create(&building)
	}

	// Create rooms (10 per building)
	var rooms []models.Room
	for _, b := range buildings {
		for i := 1; i <= 10; i++ {
			room := models.Room{
				BuildingID:       b.BuildingID,
				RoomNumber:       uint(100 + i),
				NumBedroom:       int8(rand.Intn(2) + 1), // 1–2 bedrooms
				PrivateBathrooms: int8(rand.Intn(2)),     // 0 or 1
				HasKitchen:       rand.Intn(2) == 1,
			}
			rooms = append(rooms, room)
			db.Create(&room)
		}
	}

	// Create 100 students with random preferences
	var students []models.Student
	for i := 0; i < 100; i++ {
		student := models.Student{
			WantsAC:             rand.Intn(2) == 1,
			WantsDining:         rand.Intn(2) == 1,
			WantsKitchen:        rand.Intn(2) == 1,
			WantPrivateBathroom: rand.Intn(2) == 1,
		}
		db.Create(&student)
		students = append(students, student)
	}

	// Assign 50 students to rooms (1 per room)
	var assignments []models.Assignment
	assignedCount := 0
	for _, room := range rooms {
		if assignedCount >= 50 {
			break
		}
		assignment := models.Assignment{
			StudentID:  students[assignedCount].StudentID,
			BuildingID: room.BuildingID,
			RoomNumber: room.RoomNumber,
		}
		assignments = append(assignments, assignment)
		db.Create(&assignment)
		assignedCount++
	}

	fmt.Println("University of Arkansas dorm data seeded successfully.")
}
