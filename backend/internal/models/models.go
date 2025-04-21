package models

type Assignment struct {
	StudentID  uint `gorm:"primaryKey;column:student_id" json:"student_id"`
	BuildingID uint `gorm:"column:building_id" json:"building_id"`
	RoomNumber uint `gorm:"column:room_number" json:"room_number"`
}

func (Assignment) TableName() string {
	return "assignments"
}

type Building struct {
	BuildingID   uint   `gorm:"primaryKey;autoIncrement;column:building_id" json:"building_id"`
	BuildingName string `gorm:"column:building_name;size:128" json:"building_name"`
	HasAC        bool   `gorm:"column:has_ac" json:"has_ac"`
	HasDining    bool   `gorm:"column:has_dining" json:"has_dining"`
}

func (Building) TableName() string {
	return "buildings"
}

type Room struct {
	BuildingID       uint `gorm:"column:building_id;not null;uniqueIndex:idx_building_room" json:"building_id"`
	RoomNumber       uint `gorm:"column:room_number;not null;uniqueIndex:idx_building_room" json:"room_number"`
	NumBedroom       int8 `gorm:"column:num_bedroom" json:"num_bedroom"`
	PrivateBathrooms int8 `gorm:"column:private_bathrooms" json:"private_bathrooms"`
	HasKitchen       bool `gorm:"column:has_kitchen" json:"has_kitchen"`

	Building Building `gorm:"foreignKey:BuildingID;references:BuildingID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (Room) TableName() string {
	return "rooms"
}

type Student struct {
	StudentID           uint `gorm:"primaryKey;autoIncrement;column:student_id" json:"student_id"`
	WantsAC             bool `gorm:"column:wants_ac" json:"wants_ac"`
	WantsDining         bool `gorm:"column:wants_dining" json:"wants_dining"`
	WantsKitchen        bool `gorm:"column:wants_kitchen" json:"wants_kitchen"`
	WantPrivateBathroom bool `gorm:"column:wants_private_bath" json:"wants_private_bath"`
}

func (Student) TableName() string {
	return "students"
}
