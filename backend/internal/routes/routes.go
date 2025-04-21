package routes

import (
	"github.com/eli-bosch/DBMS-final/internal/controller"

	"github.com/gorilla/mux"
)

var DBMSRoutes = func(router *mux.Router) {
	//User Routes - Add protections and logging
	router.HandleFunc("/api/student", controller.CreateStudent).Methods("POST")                               //1) Add Student to the Student Table
	router.HandleFunc("/api/assignment", controller.CreateAssignment).Methods("POST")                         //2) Add Assignment
	router.HandleFunc("/api/assignment/{building_name}", controller.FindAssignmentsByBuilding).Methods("GET") //3) Assignment in building
	router.HandleFunc("/api/rooms/{building_name}", controller.FindRoomsByBuilding).Methods("GET")            //4) All Rooms sorted by building Id
	router.HandleFunc("/api/rooms/{student_name}", controller.FindAllRoomsByPreference).Methods("GET")        //5) View all rooms by student preference
	router.HandleFunc("/api/student/{student_name}", controller.FindRoomateByPreference).Methods("GET")       //6) View all students that could room with student
	router.HandleFunc("/api/building/report", controller.ViewRoomReport).Methods("GET")                       //7) View report for each building, the number of total rooms...
	router.HandleFunc("/api/student/{student_name}", controller.FindRoomateWithAssignment).Methods("GET")     //8) View all students w/assignements that student could room with
}
