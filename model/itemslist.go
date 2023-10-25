package model

import "time"

type Item struct {
	ItemList    []ItemList `json:"itemList"`
	DateCreated int64      `json:"dateCreated"`
}
type ItemList struct {
	ItemCode string `json:"itemCode"`
	ItemDesc string `json:"itemDesc"`
	OrderQty int    `json:"orderQty"`
	Scanitem int    `json:"scanItem"`
}

type FormInputs struct {
	Invoice string `json:"invoice"`
	Scan    string `json:"scan"`
}

type Header struct {
	ID            int    `json:"id" gorm:"primaryKey"`
	InvoiceNumber string `json:"invoice"`
	BranchCode    string `json:"branchCode"`
	BmsId         string `json:"bmsId"`

	//ill change this to name --------------------------------------------------------------------
	UserId      string    `json:"userId"`
	DateCreated time.Time `json:"dateCreated"`
}

type Details struct {
	HeaderId int    `json:"header_id"`
	Header   Header `gorm:"foreignKey:HeaderId"`
	QrCode   string `json:"qrCode"`
}
