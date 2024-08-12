package models

import (
	"gorm.io/gorm"
)

type Submissions struct {
	gorm.Model
	Image   string `json:"image"`
	Caption string `json:"caption"`
	XCoordinate int `json:"x"`
	YCoordinate int `json:"y"`
	Color string `json:"color"`
}
