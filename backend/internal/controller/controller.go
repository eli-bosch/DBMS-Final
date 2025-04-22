package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/eli-bosch/DBMS-final/config"
	"github.com/eli-bosch/DBMS-final/internal/models"
	utils "github.com/eli-bosch/DBMS-final/internal/util"
	"github.com/gorilla/mux"
)


func CreateStudent(w http.ResponseWriter, r *http.Request) {
    student := &models.Student{}
    utils.ParseBody(r, student)

    db := config.GetDB()
    // GORM’s Create will INSERT and fill student.StudentID
    if err := db.Create(student).Error; err != nil {
        log.Printf("Error creating student: %v", err)
        http.Error(w, "Could not create student", http.StatusBadRequest)
        return
    }

    // student.StudentID is now set to the auto‑increment value
    res, err := json.Marshal(student)
    if err != nil {
        log.Printf("Error marshalling student: %v", err)
        http.Error(w, "Server error", http.StatusInternalServerError)
        return
    }

	fmt.Println("Student added to database")

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(res)
}



func CreateAssignment(w http.ResponseWriter, r *http.Request) {
    assignment := &models.Assignment{}
    utils.ParseBody(r, assignment)

    db := config.GetDB()

    // 1) fetch student
    var student models.Student
    if err := db.First(&student, "student_id = ?", assignment.StudentID).Error; err != nil {
        http.Error(w, "Student not found", http.StatusBadRequest)
        return
    }

    // 2) fetch building
    var building models.Building
    if err := db.First(&building, "building_id = ?", assignment.BuildingID).Error; err != nil {
        http.Error(w, "Building not found", http.StatusBadRequest)
        return
    }

    // 3) fetch room
    var room models.Room
    if err := db.
        Where("building_id = ? AND room_number = ?", assignment.BuildingID, assignment.RoomNumber).
        First(&room).Error; err != nil {
        http.Error(w, "Room not found", http.StatusBadRequest)
        return
    }

    // 4) enforce preferences
    if student.WantsAC && !building.HasAC {
        http.Error(w, "Building does not have AC", http.StatusBadRequest)
        return
    }
    if student.WantsDining && !building.HasDining {
        http.Error(w, "Building does not have Dining", http.StatusBadRequest)
        return
    }
    if student.WantsKitchen && !room.HasKitchen {
        http.Error(w, "Room does not have Kitchen", http.StatusBadRequest)
        return
    }
    if student.WantsPrivateBathroom && room.PrivateBathrooms < 1 {
        http.Error(w, "Room does not have a private bathroom", http.StatusBadRequest)
        return
    }

    // 5) create assignment
    if err := db.Create(assignment).Error; err != nil {
        log.Printf("Error creating assignment: %v", err)
        http.Error(w, "Could not create assignment", http.StatusInternalServerError)
        return
    }

	fmt.Println("Assignment created")

    // 6) return JSON
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(assignment)
}


// FindAssignmentsByBuilding returns this struct to the frontend
type AssignmentView struct {
    StudentID   uint   `json:"student_id"`
    StudentName string `json:"student_name"`
    BuildingID  uint   `json:"building_id"`
    RoomNumber  uint   `json:"room_number"`
}

func FindAssignmentsByBuilding(w http.ResponseWriter, r *http.Request) {
    buildingName := mux.Vars(r)["building_name"]

    db := config.GetDB()
    var list []AssignmentView

    // Join students & buildings, order by student name
    sql := `
    SELECT
      a.student_id,
      s.name   AS student_name,
      a.building_id,
      a.room_number
    FROM assignments a
    JOIN students s ON a.student_id = s.student_id
    JOIN buildings b ON a.building_id = b.building_id
    WHERE b.building_name = ?
    ORDER BY s.name ASC
    `

    if err := db.Raw(sql, buildingName).Scan(&list).Error; err != nil {
        log.Printf("Error querying assignments: %v", err)
        http.Error(w, "Could not fetch assignments", http.StatusBadRequest)
        return
    }

	fmt.Println("Queried list of assignments in buildling " + buildingName)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(list)
}

// FindRoomsByBuilding returns this to frontend 
type RoomView struct {
    BuildingID  uint  `json:"building_id"`
    RoomNumber  uint  `json:"room_number"`
    NumBedrooms int8  `json:"num_bedrooms"`
}

func FindAllRooms(w http.ResponseWriter, r *http.Request) {
    db := config.GetDB()

    var rooms []RoomView
    sql := `
    SELECT
      r.building_id,
      r.room_number,
      r.num_bedroom    AS num_bedrooms
    FROM rooms r
    ORDER BY r.building_id ASC, r.room_number ASC;
    `
    if err := db.Raw(sql).Scan(&rooms).Error; err != nil {
        log.Printf("Error fetching rooms: %v\n", err)
        http.Error(w, "Could not fetch rooms", http.StatusInternalServerError)
        return
    }

	fmt.Println("Loaded all rooms")

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(rooms); err != nil {
        log.Printf("Error encoding rooms JSON: %v\n", err)
        http.Error(w, "Server error", http.StatusInternalServerError)
    }
}

// FinaAllRoomsByPreference returns this to frontend
type RoomPrefView struct {
    BuildingID       uint  `json:"building_id"`
    RoomNumber       uint  `json:"room_number"`
    NumBedroom       int8  `json:"num_bedroom"`
    PrivateBathrooms int8  `json:"private_bathrooms"`
    HasKitchen       bool  `json:"has_kitchen"`
}

func FindAllRoomsByPreference(w http.ResponseWriter, r *http.Request) {
    idStr := mux.Vars(r)["student_id"]
    studentID, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
        http.Error(w, "Invalid student ID", http.StatusBadRequest)
        return
    }

    sql := `
    SELECT
      r.building_id,
      r.room_number,
      r.num_bedroom,
      r.private_bathrooms,
      r.has_kitchen
    FROM students s
    JOIN rooms r         ON 1=1
    JOIN buildings b     ON r.building_id = b.building_id
    WHERE s.student_id = ?
      AND s.wants_ac            = b.has_ac
      AND s.wants_dining        = b.has_dining
      AND s.wants_kitchen       = r.has_kitchen
      AND (
          (s.wants_private_bath = true  AND r.private_bathrooms > 0)
       OR (s.wants_private_bath = false AND r.private_bathrooms = 0)
      )
    ORDER BY r.building_id, r.room_number;
    `

    var rooms []RoomPrefView
    if err := config.GetDB().Raw(sql, studentID).Scan(&rooms).Error; err != nil {
        log.Printf("Error finding rooms: %v", err)
        http.Error(w, "Could not fetch rooms", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(rooms)
}

// struct FindRoomateByPreference returns to frontend
type RoommateView struct {
    StudentID           uint   `json:"student_id"`
    Name                string `json:"name"`
    WantsAC             bool   `json:"wants_ac"`
    WantsDining         bool   `json:"wants_dining"`
    WantsKitchen        bool   `json:"wants_kitchen"`
    WantsPrivateBathroom bool  `json:"wants_private_bath"`
}

func FindRoomateByPreference(w http.ResponseWriter, r *http.Request) {
    idStr := mux.Vars(r)["student_id"]
    studentID, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
        http.Error(w, "Invalid student ID", http.StatusBadRequest)
        return
    }

    sql := `
    SELECT 
      s2.student_id,
      s2.name,
      s2.wants_ac,
      s2.wants_dining,
      s2.wants_kitchen,
      s2.wants_private_bath
    FROM students s1
    JOIN students s2 
      ON s1.student_id != s2.student_id
     AND s1.wants_ac             = s2.wants_ac
     AND s1.wants_dining         = s2.wants_dining
     AND s1.wants_kitchen        = s2.wants_kitchen
     AND s1.wants_private_bath   = s2.wants_private_bath
    WHERE s1.student_id = ?
    ORDER BY s2.name ASC;
    `

    var mates []RoommateView
    if err := config.GetDB().Raw(sql, studentID).Scan(&mates).Error; err != nil {
        log.Printf("Error fetching roommates: %v\n", err)
        http.Error(w, "Could not fetch roommates", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(mates)
}

// Send an array of these structs to the frontend from ViewRoomReport
type ReportRow struct {
    BuildingName          string `json:"building_name"`
    TotalRooms            *int   `json:"total_rooms"`
    TotalBedrooms         *int   `json:"total_bedrooms"`
    RoomsWithAvailability *int   `json:"rooms_with_availability"`
    AvailableBedrooms     int    `json:"available_bedrooms"`
}

func ViewRoomReport(w http.ResponseWriter, r *http.Request) {
    db := config.GetDB()
    var rows []ReportRow

    sql := `
    SELECT 
      b.building_name,
      COUNT(*)                                 AS total_rooms,
      SUM(r.num_bedroom)                       AS total_bedrooms,
      SUM(CASE WHEN COALESCE(a.assigned_count,0) < r.num_bedroom THEN 1 ELSE 0 END)  AS rooms_with_availability,
      SUM(CASE WHEN COALESCE(a.assigned_count,0) < r.num_bedroom THEN r.num_bedroom-COALESCE(a.assigned_count,0) ELSE 0 END) AS available_bedrooms
    FROM rooms r
    JOIN buildings b 
      ON r.building_id = b.building_id
    LEFT JOIN (
      SELECT building_id, room_number, COUNT(*) AS assigned_count
      FROM assignments
      GROUP BY building_id, room_number
    ) a 
      ON r.building_id = a.building_id
     AND r.room_number = a.room_number
    GROUP BY b.building_name

    UNION ALL

    SELECT
      'TOTAL'                                         AS building_name,
      SUM(t.total_rooms)           AS total_rooms,
      SUM(t.total_bedrooms)        AS total_bedrooms,
      SUM(t.rooms_with_availability) AS rooms_with_availability,
      SUM(t.available_bedrooms)    AS available_bedrooms
    FROM (
      SELECT 
        COUNT(*)                                 AS total_rooms,
        SUM(r.num_bedroom)                       AS total_bedrooms,
        SUM(CASE WHEN COALESCE(a.assigned_count,0) < r.num_bedroom THEN 1 ELSE 0 END)  AS rooms_with_availability,
        SUM(CASE WHEN COALESCE(a.assigned_count,0) < r.num_bedroom THEN r.num_bedroom-COALESCE(a.assigned_count,0) ELSE 0 END) AS available_bedrooms
      FROM rooms r
      LEFT JOIN (
        SELECT building_id, room_number, COUNT(*) AS assigned_count
        FROM assignments
        GROUP BY building_id, room_number
      ) a 
        ON r.building_id = a.building_id
       AND r.room_number = a.room_number
    ) AS t

    ORDER BY (building_name = 'TOTAL'), building_name;
    `

    if err := db.Raw(sql).Scan(&rows).Error; err != nil {
        log.Printf("Error fetching report: %v\n", err)
        http.Error(w, "Could not fetch report", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(rows)
}