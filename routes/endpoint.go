package routes

import (
	"couchbase/model"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var (
	store         = session.New()
	ItemsResponse []model.ItemList
	itemCodes     []string
)

func GetItemList(c *fiber.Ctx) error {
	var formInput model.FormInputs
	c.BodyParser(&formInput)

	if CheckInvoiceIfExist(formInput.Invoice) {
		return c.Render("index", fiber.Map{
			"UserName": GetSessionValue(c, "userId"),
			"Success":  false,
			"Message":  "Invoice Number Already Encoded!",
		})
	}

	Items := GetList(formInput.Invoice)
	if len(Items) == 0 {
		return c.Render("index", fiber.Map{
			"UserName": GetSessionValue(c, "userId"),
			"Success":  false,
			"Message":  "Invalid Invoice Number!",
		})
	}

	// if IsExpired(Items[0].DateCreated) {
	// 	return c.Render("index", fiber.Map{
	// 		"Success": false,
	// 		"Message": "Invoice not Found",
	// 	})
	// }

	msgResp := CheckProductIfIncludedInList(Items[0].ItemList)
	if msgResp != "" {
		return c.Render("index", fiber.Map{
			"UserName": GetSessionValue(c, "userId"),
			"Success":  false,
			"Message":  msgResp,
		})
	}

	ItemList := ItemsResponse
	invoiceNumber := SavingSession(c, "invoice", formInput.Invoice)
	itemCodes = []string{}
	return c.Render("scan", fiber.Map{
		"UserName": GetSessionValue(c, "userId"),
		"Invoice":  invoiceNumber,
		"Data":     ItemList,
	})
}

func CancelScan(c *fiber.Ctx) error {
	ItemsResponse = []model.ItemList{}
	return c.Redirect("/?" + "branchCode=" + GetSessionValue(c, "branchCode") + "&bmsId=" + GetSessionValue(c, "bmsId") + "&userName=" + GetSessionValue(c, "userId"))
}

func QrScan(c *fiber.Ctx) error {
	var formInput model.FormInputs
	c.BodyParser(&formInput)

	itemCode := TrimAndGetFirst6Characters(formInput.Scan)
	itemList := ItemsResponse
	invoiceNumber := GetSessionValue(c, "invoice")

	for i, v := range itemList {

		if itemCode == v.ItemCode {
			if len(formInput.Scan) != 29 {

				return c.Render("scan", fiber.Map{
					"UserName": GetSessionValue(c, "userId"),
					"Success":  false,
					"Message":  "Illegal Qr Code Format!",
					"Data":     itemList,
					"Invoice":  invoiceNumber,
				})
			}

			if v.OrderQty > 0 {
				AppendingToitemCodeList(formInput.Scan, &itemCodes)
				ItemsResponse[i].OrderQty--
				ItemsResponse[i].Scanitem++

				isDone, headerId := IsDone(itemList, c)
				if isDone {
					InsertToDetails(itemCodes, headerId)
					return c.Render("index", &fiber.Map{
						"UserName": GetSessionValue(c, "userId"),
						"Success":  true,
						"Message":  "Transaction Complete!",
					})
				}
				return c.Render("scan", fiber.Map{
					"UserName": GetSessionValue(c, "userId"),
					"Data":     itemList,
					"Invoice":  invoiceNumber,
				})
			}

			return c.Render("scan", fiber.Map{
				"UserName": GetSessionValue(c, "userId"),
				"Success":  false,
				"Message":  "Invoice Item Already Counted In Full!",
				"Data":     itemList,
				"Invoice":  invoiceNumber,
			})

		}
	}

	return c.Render("scan", fiber.Map{
		"UserName": GetSessionValue(c, "userId"),
		"Success":  false,
		"Message":  "Qr Code Mismatch!",
		"Data":     itemList,
		"Invoice":  invoiceNumber,
	})
}

func AppendingToitemCodeList(itemCode string, itemCodeList *[]string) {
	fmt.Println("scanned code:", itemCode)
	*itemCodeList = append(*itemCodeList, itemCode)

}
