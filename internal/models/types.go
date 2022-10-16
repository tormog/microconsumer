package models

type DataStore struct {
	ID   string      `json:"id"`
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}
