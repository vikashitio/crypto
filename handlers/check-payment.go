package handlers

import (
	"fmt"
	"os"
	"template/database"
	"template/function"
	"template/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Function for Display invoice payment page with qr code
func PaymentQRView(c *fiber.Ctx) error {

	PaymentID := c.Query("iid")
	//fmt.Println(PaymentID)
	//=============Fetch Invoice Details by trackid===============
	invoiceData := models.Invoice_Master{}
	database.DB.Db.Table("invoice").Where("trackid = ?", PaymentID).Find(&invoiceData)
	//=============Fetch coin list ===============
	coinList := []models.CoinList{}
	database.DB.Db.Table("coin_list").Order("coin ASC").Where("status = ?", 1).Find(&coinList)
	//fmt.Println(coinList)
	var commonURL = os.Getenv("CommonURL")
	return c.Render("checkout-pay", fiber.Map{
		"Title":       "Checkout",
		"Subtitle":    "Checkout",
		"CoinList":    coinList,
		"InvoiceData": invoiceData,
		"CommonURL":   commonURL,
	})
}

// Get Coin Balance
type FailedStatus struct {
	Status             string    `json:"status,omitempty"`
	Response_timestamp time.Time `json:"response_timestamp,omitempty"`
}

// Function for Display Failed payment page
func FailedView(c *fiber.Ctx) error {

	transID := c.Params("TransID")
	fmt.Println(transID)

	//=============Update Transaction status from Transaction id===============
	currentTime := time.Now()
	database.DB.Db.Table("transaction").Where("transaction_id = ?", transID).UpdateColumns(FailedStatus{Status: "FAILED", Response_timestamp: currentTime})

	//=============Fetch Transaction details from Transaction id===============
	transData := models.Transaction_Pay{}
	database.DB.Db.Table("transaction").Where("transaction_id = ?", transID).Find(&transData)
	mID := transData.Client_id

	//  Email///
	template_code := "PAYMENT-STATUS"
	getName := function.GetNameByMID(mID)
	getEmail := function.GetEmailByMID(mID)
	hashCode := transData.Response_hash
	num := transData.Receivedamount
	amount := fmt.Sprintf("%.2f", num) // Convert to string with 2 decimal places
	crypto := transData.Receivedcurrency
	status := transData.Status

	emailData := models.EmailData{FullName: getName, Email: getEmail, Status: status, Amount: amount, Crypto: crypto, HashCode: hashCode, TransID: transID}

	err := function.SendEmail(template_code, emailData)
	if err != nil {
		fmt.Println("issue sending verification email")
	}

	return c.Render("failed", fiber.Map{
		"Title":     "Failed Payment",
		"Subtitle":  "Failed Payment",
		"TransID":   transID,
		"TransData": transData,
	})
}

// Function for Display Success payment page
func SuccessView(c *fiber.Ctx) error {

	transID := c.Params("TransID")
	fmt.Println(transID)
	//=============Fetch Invoice Details by trackid===============

	transData := models.Transaction_Pay{}
	database.DB.Db.Table("transaction").Where("transaction_id = ?", transID).Find(&transData)
	mID := transData.Client_id

	//  Email///
	template_code := "PAYMENT-STATUS"
	getName := function.GetNameByMID(mID)
	getEmail := function.GetEmailByMID(mID)
	hashCode := transData.Response_hash
	num := transData.Receivedamount
	amount := fmt.Sprintf("%.2f", num) // Convert to string with 2 decimal places
	crypto := transData.Receivedcurrency
	status := transData.Status

	emailData := models.EmailData{FullName: getName, Email: getEmail, Status: status, Amount: amount, Crypto: crypto, HashCode: hashCode, TransID: transID}

	err := function.SendEmail(template_code, emailData)
	if err != nil {
		fmt.Println("issue sending verification email")
	}

	return c.Render("success", fiber.Map{
		"Title":     "Success Payment",
		"Subtitle":  "Success Payment",
		"TransID":   transID,
		"TransData": transData,
	})
}

// Function for Display  payment detail on checkout
func PayQRView(c *fiber.Ctx) error {

	PaymentID := c.Query("iid")
	//fmt.Println(PaymentID)
	//=============Fetch Invoice Details by trackid===============
	invoiceData := models.Invoice_Master{}
	database.DB.Db.Table("invoice").Where("trackid = ?", PaymentID).Find(&invoiceData)
	//=============Fetch coin list ===============
	coinList := []models.CoinList{}
	database.DB.Db.Table("coin_list").Order("coin ASC").Where("status = ?", 1).Find(&coinList)
	fmt.Println(invoiceData)
	var commonURL = os.Getenv("CommonURL")
	return c.Render("checkout-pay-views", fiber.Map{
		"Title":       "Checkout",
		"Subtitle":    "Checkout",
		"CoinList":    coinList,
		"InvoiceData": invoiceData,
		"CommonURL":   commonURL,
	})
}
