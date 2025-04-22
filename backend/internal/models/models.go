// internal/models/models.go
package models

// Assignment represents a student‑room assignment.
// StudentID is the primary key.
type Assignment struct {
    StudentID  uint `gorm:"primary_key;column:student_id" json:"student_id"`
    BuildingID uint `gorm:"column:building_id"       json:"building_id"`
    RoomNumber uint `gorm:"column:room_number"       json:"room_number"`
}

func (Assignment) TableName() string {
    return "assignments"
}

// Building represents a dorm building.
type Building struct {
    BuildingID   uint   `gorm:"primary_key;AUTO_INCREMENT;column:building_id" json:"building_id"`
    BuildingName string `gorm:"column:building_name;size:128"               json:"building_name"`
    HasAC        bool   `gorm:"column:has_ac"                                json:"has_ac"`
    HasDining    bool   `gorm:"column:has_dining"                            json:"has_dining"`
}

func (Building) TableName() string {
    return "buildings"
}

// Room represents a room in a building.
// A composite unique index ensures (building_id, room_number) is unique.
type Room struct {
    BuildingID       uint     `gorm:"column:building_id;not null;unique_index:idx_building_room" json:"building_id"`
    RoomNumber       uint     `gorm:"column:room_number;not null;unique_index:idx_building_room" json:"room_number"`
    NumBedroom       int8     `gorm:"column:num_bedroom"                                  json:"num_bedroom"`
    PrivateBathrooms int8     `gorm:"column:private_bathrooms"                            json:"private_bathrooms"`
    HasKitchen       bool     `gorm:"column:has_kitchen"                                  json:"has_kitchen"`
    Building         Building `gorm:"foreignkey:BuildingID;association_foreignkey:BuildingID"`
}

func (Room) TableName() string {
    return "rooms"
}

// Student represents a freshman looking for a roommate.
type Student struct {
    StudentID           uint `gorm:"primary_key;AUTO_INCREMENT;column:student_id" json:"student_id"`
	Name                string `gorm:"column:name;size:128;not null"               json:"name"`
    WantsAC             bool `gorm:"column:wants_ac"                           json:"wants_ac"`
    WantsDining         bool `gorm:"column:wants_dining"                       json:"wants_dining"`
    WantsKitchen        bool `gorm:"column:wants_kitchen"                      json:"wants_kitchen"`
    WantsPrivateBathroom bool `gorm:"column:wants_private_bath"                 json:"wants_private_bath"`
}

func (Student) TableName() string {
    return "students"
}
