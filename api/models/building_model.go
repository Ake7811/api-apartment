package models

import (
	"html"
	"strings"
	"time"
)

type Building struct {
	Id              uint64     `gorm:"primay_key;auto_increment" json:"id"`
	Building_name   string     `gorm:"size:200;not null" json:"building_name"`
	Remark          string     `gorm:"size:255;" json:"remark"`
	Building_status string     `gorm:"size:1;not null" json:"building_status"` // (0 = delete , 1 = active , 2 = inactive)
	Create_by       uint64     `gorm:"not null" json:"create_by"`
	Create_At       time.Time  `gorm:"default:current_timestamp()" json:"create_at"`
	Update_by       uint64     `json:"update_by"`
	Update_At       *time.Time `gorm:"default:current_timestamp()" json:"update_at"`
}

func (u *Building) BeforeSave() error {
	return nil
}

func (u *Building) Prepare() {
	// For create
	u.Id = 0
	u.Building_name = html.EscapeString(strings.TrimSpace(u.Building_name))
	u.Remark = html.EscapeString(strings.TrimSpace(u.Remark))
	u.Building_status = "1" // Active
	u.Create_At = time.Now()
	u.Update_by = 0
	u.Update_At = nil
}

func (u *Building) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
	case "delete":
	default:
	}
	return nil
}
