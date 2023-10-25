package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {

	// app.Group("searialization/v1")

	app.Get("/", func(c *fiber.Ctx) error {
		queryParam := c.Queries()

		//ill put it in the session
		branchCode := queryParam["branchCode"]
		bmsId := queryParam["bmsId"]
		userId := queryParam["userName"]

		//--------------------------ill add this later------------------------------//
		// headerValue := c.Get("Authorization")
		// if headerValue == "" || headerValue != "my set token" {
		// 	return c.Render("errorpage", fiber.Map{
		// 		"Message": "Illegal Request",
		// 	})
		// }
		//--------------------------------------------------------//

		//saving params value in session
		SavingSession(c, "branchCode", branchCode)
		SavingSession(c, "bmsId", bmsId)
		SavingSession(c, "userId", uppercaseFirstLetter(userId))

		fmt.Printf("branchCode:%v \n bmsId:%s \n userId:%s", GetSessionValue(c, "branchCode"), GetSessionValue(c, "bmsId"), userId)

		return c.Render("index", &fiber.Map{
			"UserName": GetSessionValue(c, "userId"),
		})
	})
	// app.Get("/scan", func(c *fiber.Ctx) error {
	// 	return c.Render("scan", nil)
	// })
	app.Post("/getItemList", GetItemList)
	app.Get("/cancel", CancelScan)
	app.Post("/scanQR", QrScan)
}
