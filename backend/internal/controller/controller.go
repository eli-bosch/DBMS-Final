package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/eli-bosch/DBMS-final/config"
	"github.com/eli-bosch/DBMS-final/internal/models"
	utils "github.com/eli-bosch/DBMS-final/internal/util"
	"github.com/gorilla/mux"
)

func CreateStudent(w http.ResponseWriter, r *http.Request) {
	student := &models.Student{} //Parses JSON POST request
	utils.ParseBody(r, student)

	db := config.GetDB() //Creates DB entry
	result := db.Exec(`
	INSERT INTO students (wants_ac, wants_dining, wants_kitchen, wants_private_bath)
	VALUES (?, ?, ?, ?, ?)
	`, student.WantsAC, student.WantsDining, student.WantsKitchen, student.WantPrivateBathroom)

	if result.Error != nil {
		fmt.Printf("Error creating student: %v", result.Error)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var id uint
	db.Raw("SELECT LAST_INSERT_ID()").Scan(&id)
	student.StudentID = id

	res, err := json.Marshal(student)
	if err != nil {
		fmt.Println("Error marshalling student")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json") //Sends body back to requester
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func CreateAssignment(w http.ResponseWriter, r *http.Request) {
	assignment := &models.Assignment{} //Parses JSON POST request
	utils.ParseBody(r, assignment)

	db := config.GetDB() //Creates DB entry
	result := db.Exec(`INSERT INTO assignments (student_id, building_id, room_number) VALUES
	(?,?,?)`, assignment.StudentID, assignment.BuildingID, assignment.RoomNumber)

	if result.Error != nil {
		fmt.Printf("Error creating assignment: %v", result.Error)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := json.Marshal(assignment) //Marshalls json body
	if err != nil {
		fmt.Println("Error marshalling assignment")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json") //Sends body back to requester
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func FindAssignmentsByBuilding(w http.ResponseWriter, r *http.Request) {
	var assignments []models.Assignment

	vars := mux.Vars(r)
	name := vars["building_name"]

	db := config.GetDB()
	err := db.Raw(`
	SELECT a.student_id, a.room_number, a.building_id, s.name AS student_name
	FROM assignments a
	JOIN students s ON a.student_id = s.student_id
	JOIN building b ON a.building_id = b.building_id
	WHERE b.building_name = ?
	ORDER BY s.name ASC`, name).Scan(&assignments).Error

	if err != nil {
		fmt.Printf("Error querying assignments: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := json.Marshal(&assignments)
	if err != nil {
		fmt.Printf("Error marshalling assignments: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func FindRoomsByBuilding(w http.ResponseWriter, r *http.Request) {
	var rooms []models.Room
	vars := mux.Vars(r)
	name := vars["building_name"]

	db := config.GetDB()
	err := db.Raw(`
	SELECT r.*
	FROM rooms r
	JOIN buildings b ON r.buildin_id = b.building_id
	WHERE b.name = ?`, name).Scan(&rooms).Error

	if err != nil {
		fmt.Printf("Error finding rooms: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := json.Marshal(&rooms)
	if err != nil {
		fmt.Println("Error marshalling rooms")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func FindAllRoomsByPreference(w http.ResponseWriter, r *http.Request) {
	var rooms []models.Room
	vars := mux.Vars(r)
	name := vars["student_name"]

	db := config.GetDB()
	err := db.Raw(`
	SELECT r.*
	FROM students s
	JOIN rooms r ON 1=1
	JOIN buildings b ON r.building_id = b.building_id
	WHERE s.name = ?
		AND (s.wants_ac = false OR b.has_ac = true)
		AND (s.wants_kitchen = false OR b.has_kitchen = true)
		AND (s.wants_private_bathroom = false OR r.private_bathroom > 0)
		AnD (s.wants_dining = false OR b.has_dining = true)
	`, name).Scan(&rooms).Error

	if err != nil {
		fmt.Printf("Error finding rooms: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := json.Marshal(&rooms)
	if err != nil {
		fmt.Println("Error while marshalling rooms")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func FindRoomateByPreference(w http.ResponseWriter, r *http.Request) {
	var roomates []models.Student
	vars := mux.Vars(r)
	name := vars["student_name"]

	db := config.GetDB()
	err := db.Raw(`
	SELECT s2.*
	FROM students s1
	JOIN students s2 ON
		s1.name = ?
		AND s1.students_id != s2.student_id
		AND s1.wants_ac = s2.wants_acnil
		AND s1.wants_dining = s2.wants_dining
		AND s1.wants_kitchen = s2.wants_kitchen
		AND s1.wants_private_bath = s2.wants_private_bath
	`, name).Scan(&roomates).Error

	if err != nil {
		fmt.Printf("Error fetching roomates: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := json.Marshal(&roomates)
	if err != nil {
		fmt.Println("Error marshalling roomates")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func ViewRoomReport(w http.ResponseWriter, r *http.Request) {
	type Report struct {
		BuildingName     string `json:"building_name"`
		TotalRooms       *int   `json:"total_rooms"`
		TotalBedrooms    *int   `json:"rooms_with_availability"`
		AvaiableBedrooms int    `json:"available_bedrooms"`
	}

	var report []Report
	db := config.GetDB()

	err := db.Raw(`
	WITH room_assignments AS (
    SELECT
        r.building_id,
        r.room_number,
        b.name AS building_name,
        r.num_bedroom,
        COUNT(a.student_id) AS assigned_count
    FROM rooms r
    JOIN buildings b ON r.building_id = b.building_id
    LEFT JOIN assignments a ON r.building_id = a.building_id AND r.room_number = a.room_number
    GROUP BY r.building_id, r.room_number, b.name, r.num_bedroom
	),
	building_summary AS (
    SELECT
        building_name,
        COUNT(*) AS total_rooms,
        SUM(num_bedroom) AS total_bedrooms,
        COUNT(CASE WHEN assigned_count < num_bedroom THEN 1 END) AS rooms_with_availability,
        SUM(CASE WHEN assigned_count < num_bedroom THEN num_bedroom - assigned_count ELSE 0 END) AS available_bedrooms
    FROM room_assignments
    GROUP BY building_name
	),
	campus_summary AS (
    SELECT
        'TOTAL' AS building_name,
        NULL AS total_rooms,
        NULL AS total_bedrooms,
        NULL AS rooms_with_availability,
        SUM(available_bedrooms) AS available_bedrooms
    FROM building_summary
	)

	SELECT * FROM building_summary
	UNION ALL
	SELECT * FROM campus_summary;
	`).Scan(&report).Error

	if err != nil {
		fmt.Printf("Error fetching reports: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := json.Marshal(&report)
	if err != nil {
		fmt.Println("Error while marshalling reports")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func FindRoomateWithAssignment(w http.ResponseWriter, r *http.Request) {
}
