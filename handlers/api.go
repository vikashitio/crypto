package handlers

import (
	"fmt"
	"strconv"
	"template/database"
	"template/function"
	"template/models"

	"github.com/gofiber/fiber/v2"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type APIResponsePaylinkSuccess struct {
	Status string
	PayURL string
}
type APIResponseFailed struct {
	Status string
	Error  string
}

type APIResponseBalanceSuccess struct {
	Receivedcurrency string
	Balance          string
}

// Function for Display Payment form in merchant section
func ApiPaymentLink(c *fiber.Ctx) error {
	apiError := ""
	// Retrieve a specific header
	apikey := c.Get("Apikey")

	MID, errorx := function.GetMIDByApikey(apikey)
	// Retrieve data from header
	price_currency := c.FormValue("Currency")
	price_amount := c.FormValue("Amount")
	requestedamount, err := strconv.ParseFloat(price_amount, 64)
	if err != nil {
		fmt.Println(err)
	}

	productName := c.FormValue("ProductName")
	description := c.FormValue("Description")

	if errorx != "" {
		fmt.Println(errorx)
		apiError = errorx
	} else if price_currency == "" {
		apiError = "Currency Not Found"
	} else if price_amount == "" {
		apiError = "Amount Not Found"
	} else if productName == "" {
		apiError = "Product Name Name Not Found"
	} else if description == "" {
		apiError = "Description Not Found"
	} else {

		// Generate randomly Transaction ID
		trackID, err := function.GenerateRandomID(16) // 16 bytes will give us a 32 character hex string
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to generate random ID",
			})
		}

		Ip := c.Context().RemoteIP().String()
		qry := models.Invoice_Master{Client_id: MID, Requestedamount: requestedamount, Requestedcurrency: price_currency, Product_name: productName, Description: description, Ip: Ip, Trackid: trackID}
		result := database.DB.Db.Table("invoice").Select("client_id", "requestedamount", "requestedcurrency", "product_name", "description", "ip", "trackid").Create(&qry)
		invoice_id := strconv.FormatUint(uint64(qry.Invoice_id), 10)
		fmt.Println(invoice_id)
		if result.Error != nil {
			fmt.Println("Data Not Inserted")

		}

		PayURL := "http://localhost:3000/pay?iid=" + trackID
		status := "Ok"

		response := APIResponsePaylinkSuccess{
			Status: status,
			PayURL: PayURL,
		}

		return c.JSON(response)
	}
	status := "Fail"
	response := APIResponseFailed{
		Status: status,
		Error:  apiError,
	}

	return c.JSON(response)
}
func ApiBalanceByCrypto(c *fiber.Ctx) error {

	apiError := ""
	// Retrieve a specific header
	apikey := c.Get("Apikey")

	MID, errorx := function.GetMIDByApikey(apikey)
	fmt.Println(MID)
	if errorx != "" {
		fmt.Println(errorx)
		apiError = errorx

	} else {

		assetList := []APIResponseBalanceSuccess{}
		var totalWallet int64
		// fetch query for wallet with balance
		database.DB.Db.Table("transaction").Select("assetid, receivedcurrency, SUM(receivedamount)  as balance").Where("client_id = ? AND status = ?", MID, "SUCCESS").Group("assetid,receivedcurrency").Having("COUNT(assetid) > ?", 0).Order("assetid ASC").Find(&assetList).Count(&totalWallet)

		return c.JSON(assetList)
	}

	status := "Fail"
	response := APIResponseFailed{
		Status: status,
		Error:  apiError,
	}

	return c.JSON(response)

}
func ApiTransactionByTransID(c *fiber.Ctx) error {

	apiError := ""
	// Retrieve a specific header
	apikey := c.Get("Apikey")
	TransID := c.Params("TransID")
	fmt.Println("TransID==>", TransID)

	MID, errorx := function.GetMIDByApikey(apikey)
	fmt.Println(MID)
	if errorx != "" {
		fmt.Println(errorx)
		apiError = errorx
	} else if TransID == "" {
		apiError = "TransID Not Found"

	} else {

		transData := models.Transaction_Pay{}
		database.DB.Db.Table("transaction").Where("transaction_id = ? AND client_id = ? ", TransID, MID).Omit("client_id", "assetid", "response_json").Find(&transData)

		return c.JSON(transData)
	}

	status := "Fail"
	response := APIResponseFailed{
		Status: status,
		Error:  apiError,
	}

	return c.JSON(response)

}
func ApiTransactionList(c *fiber.Ctx) error {

	apiError := ""
	// Retrieve a specific header
	apikey := c.Get("Apikey")

	MID, errorx := function.GetMIDByApikey(apikey)
	fmt.Println(MID)
	if errorx != "" {
		fmt.Println(errorx)
		apiError = errorx

	} else {

		transData := []models.Transaction_Pay{}
		database.DB.Db.Table("transaction").Where("client_id = ? ", MID).Omit("client_id", "assetid", "response_json").Order("id DESC").Limit(10).Find(&transData)

		return c.JSON(transData)
	}

	status := "Fail"
	response := APIResponseFailed{
		Status: status,
		Error:  apiError,
	}

	return c.JSON(response)

}
