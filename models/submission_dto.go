package models

type SubmissionCreateRequest struct {
	ID      uint   `json:"-"`
	Image   string `json:"image"`
	Caption string `json:"caption"`
	XCoordinate int `json:"x"`
	YCoordinate int `json:"y"`
	Color string `json:"color"`
}

func (SubmissionCreateRequest) TableName() string { return "submissions" }
