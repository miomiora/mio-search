package model

type SearchDTO struct {
	Size int    `json:"size"`
	Page int    `json:"page"`
	Text string `json:"text"`
	Type string `json:"type"`
}
