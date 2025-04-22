package routes

import (
	"github.com/eli-bosch/DBMS-final/internal/controller"

	"github.com/gorilla/mux"
)

var DBMSRoutes = func(router *mux.Router) {
	//User Routes - Add protections and logging
	router.HandleFunc("/api/student", controller.CreateStudent).Methods("POST", "OPTIONS")                               //1) Add Student to the Student Table
	router.HandleFunc("/api/assignment", controller.CreateAssignment).Methods("POST")                                    //2) Add Assignment
	router.HandleFunc("/api/assignment/{building_name}", controller.FindAssignmentsByBuilding).Methods("GET", "OPTIONS") //3) Assignment in building
	router.HandleFunc("/api/rooms", controller.FindAllRooms).Methods("GET")                                              //4) All Rooms sorted by building Id
	router.HandleFunc("/api/preference/{student_id}", controller.FindAllRoomsByPreference).Methods("GET")                //5) View all rooms by student preference
	router.HandleFunc("/api/student/{student_id}", controller.FindRoomateByPreference).Methods("GET")                    //6) View all students that could room with student
	router.HandleFunc("/api/building/report", controller.ViewRoomReport).Methods("GET")                                  //7) View report for each building, the number of total rooms...
}
