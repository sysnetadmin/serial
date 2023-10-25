package model

type Products struct {
	ID          int    `json:"id" gorm:"primaryKey"`
	Barcode     string `json:"barcode"`
	Item        string `json:"item"`
	Description string `json:"desc"`
}
