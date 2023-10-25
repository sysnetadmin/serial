package routes

import (
	"couchbase/database"
	"couchbase/model"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/couchbase/gocb/v2"
	"github.com/gofiber/fiber/v2"
)

func TrimAndGetFirst6Characters(input string) string {
	// Trim leading and trailing spaces
	trimmed := strings.TrimSpace(input)

	// Take the first 6 characters
	if len(trimmed) <= 6 {
		return trimmed
	}
	return trimmed[:6]
}

// checking if all the item qty == 0
func IsDone(items []model.ItemList, c *fiber.Ctx) (bool, int) {
	for _, v := range items {
		if v.OrderQty > 0 {
			return false, 0
		}
	}

	//inserting into headers  -change it when needed
	ItemsResponse = []model.ItemList{}
	header := model.Header{
		InvoiceNumber: GetSessionValue(c, "invoice"),
		BranchCode:    GetSessionValue(c, "branchCode"),
		BmsId:         GetSessionValue(c, "bmsId"),
		UserId:        GetSessionValue(c, "userId"),
		DateCreated:   time.Now(),
	}

	if err := database.DB.Debug().Create(&header).Error; err != nil {
		panic(err)
	}

	return true, header.ID
}

// check if the invoice number exist in the local db
func CheckInvoiceIfExist(invoice string) bool {
	var transRespo model.Header
	row := database.DB.Debug().Find(&transRespo, "invoice_number=?", invoice).RowsAffected
	return row > 0
}

// querying in the couchbase
func GetList(invoiceNumber string) []model.Item {
	var Items []model.Item
	query := fmt.Sprintf("select  dateCreated,itemList from pcbms where javaClass='Order' and invoiceNo='%s'", invoiceNumber)

	rows, err := database.Cluster.Query(query, &gocb.QueryOptions{})
	if err != nil {
		log.Fatalf("Error executing N1QL query: %v", err)
	}
	defer rows.Close()

	// Iterate through the query results and parse them into a struct
	for rows.Next() {
		var item model.Item

		if err := rows.Row(&item); err != nil {
			log.Fatalf("Error parsing query result: %v", err)
		}
		Items = append(Items, item)
	}

	return Items

}

func IsExpired(dateCreated int64) bool {
	toSeconds := dateCreated / 1000
	dateExpired := time.Unix(toSeconds, 0).Add(4 * 24 * time.Hour).Unix()
	currentTime := time.Now().Unix()

	return currentTime > dateExpired
}

func SavingSession(c *fiber.Ctx, key, value string) interface{} {
	sess, err := store.Get(c)
	if err != nil {
		return err
	}
	sess.Set(key, value)
	invoiceNumber := sess.Get(key)
	sess.Save()

	return invoiceNumber
}

func GetSessionValue(c *fiber.Ctx, key string) string {
	sess, err := store.Get(c)
	if err != nil {
		return ""
	}

	invoiceNumber := sess.Get(key)
	return invoiceNumber.(string)
}

func CheckProductIfIncludedInList(ItemList []model.ItemList) string {
	var productList []model.Products
	codes := []string{}

	for _, v := range ItemList {
		codes = append(codes, v.ItemCode)
	}

	database.DB.Debug().Where("item IN (?)", codes).Find(&productList)
	if len(productList) == 0 {
		return "Invoice Items Not Registered!"
	}

	ItemsResponse = []model.ItemList{}
	for _, product := range productList {
		for _, item := range ItemList {
			if product.Item == item.ItemCode {
				ItemsResponse = append(ItemsResponse, item)
			}
		}
	}

	return ""
}

func isGreaterThanOne(orderQuantity int) bool {
	return orderQuantity >= 2
}

func scanMultiple(i, scanQuantity int, itemList *[]model.ItemList, c *fiber.Ctx) error {
	invoiceNumber := GetSessionValue(c, "invoice")
	ItemsResponse[i].OrderQty -= scanQuantity
	ItemsResponse[i].Scanitem += scanQuantity
	isDone, _ := IsDone(*itemList, c)
	if isDone {
		return c.Render("index", &fiber.Map{
			"Success": true,
			"Message": "Transaction Complete!",
		})
	}
	return c.Render("scan", fiber.Map{
		"Data":    itemList,
		"Invoice": invoiceNumber,
	})
}

func scanItemOnce(i, scanQuantity int, itemList *[]model.ItemList, c *fiber.Ctx) error {
	invoiceNumber := GetSessionValue(c, "invoice")
	ItemsResponse[i].OrderQty--
	ItemsResponse[i].Scanitem++
	isDone, _ := IsDone(*itemList, c)
	if isDone {

		return c.Render("index", &fiber.Map{
			"Success": true,
			"Message": "Transaction Complete!",
		})
	}
	return c.Render("scan", fiber.Map{
		"Data":    itemList,
		"Invoice": invoiceNumber,
	})
}

func TraverseITemList(itemList *[]model.ItemList, c *fiber.Ctx, itemCode string, scanQuantity int) error {
	invoiceNumber := GetSessionValue(c, "invoice")
	scanQuantitySuccess := false

	for i, v := range *itemList {
		if itemCode == v.ItemCode && v.OrderQty > 0 {
			if scanQuantity > 0 {
				return scanMultiple(i, scanQuantity, itemList, c)
			}
			if isGreaterThanOne(v.OrderQty) {
				scanQuantitySuccess = true
				return c.Render("scan", &fiber.Map{
					"ScanQuantity": scanQuantitySuccess,
					"Item":         v.ItemDesc,
					"ItemQuan":     v.OrderQty,
					"Data":         itemList,
					"Invoice":      invoiceNumber,
					"ItemCode":     itemCode,
				})
			}
			return scanItemOnce(i, scanQuantity, itemList, c)
		}
	}
	return c.Render("scan", fiber.Map{
		"Success": false,
		"Message": "Item not in the list or already scanned",
		"Data":    itemList,
		"Invoice": invoiceNumber,
	})
}

func InsertToDetails(itemCodes []string, header_id int) {

	for _, itemCode := range itemCodes {
		var details model.Details
		database.DB.Debug().Raw("insert into details (header_id,qr_code) values(?,?)", header_id, itemCode).Find(&details)

	}
}

//balikan ko bukas

func uppercaseFirstLetter(s string) string {
	if s == "" {
		return s
	}
	output := []string{}

	for _, v := range strings.Fields(s) {

		v := strings.ToUpper(v[:1]) + v[1:]

		output = append(output, v)
	}

	// Convert the first character to uppercase
	return strings.Join(output, " ")
}
