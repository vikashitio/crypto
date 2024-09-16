package handlers

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"template/database"
	"template/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

// For Admin
// Function for Display Invoice Listing in admin section
func AdminInvoiceListView(c *fiber.Ctx) error {

	AdminSession(c)
	sess, _ := store.Get(c)
	adminData := sess.Get("AdminData")
	// Get query parameters for page and limit
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid page number")
	}
	pageLimit := "10"
	limit, err := strconv.Atoi(c.Query("limit", pageLimit))
	if err != nil || limit < 1 {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid limit number")
	}

	// Calculate offset
	offset := (page - 1) * limit

	transactionList := []models.Invoice_Master{}
	database.DB.Db.Table("invoice").Order("invoice_id desc").Limit(limit).Offset(offset).Find(&transactionList)

	var total int64
	database.DB.Db.Table("invoice").Count(&total)

	//fmt.Println(total)

	// Prepare pagination data
	totalPage := total / 10
	//fmt.Println(totalPage)
	nextPage := page + 1
	prevPage := page - 1
	if page == 1 {
		prevPage = 0
	}

	if page >= int(totalPage+1) {
		nextPage = 0
	}

	Alerts := sess.Get("AlertX")
	if Alerts != "" {
		sess.Delete("AlertX")
		if err := sess.Save(); err != nil {
			panic(err)
		}
	}
	var commonURL = os.Getenv("CommonURL")

	//fmt.Println(transactionList)
	return c.Render("admin/invoice-list", fiber.Map{
		"Title":           "Invoices",
		"Subtitle":        "Invoices",
		"AlertX":          Alerts,
		"AdminData":       adminData,
		"CommonURL":       commonURL,
		"TransactionList": transactionList,
		"NextPage":        nextPage,
		"PrevPage":        prevPage,
		"Limit":           limit,
		"Count":           total,
	})
}

// Function for Display Transaction Listing in admin section
func AdminTransactionsView(c *fiber.Ctx) error {

	AdminSession(c)
	sess, _ := store.Get(c)
	adminData := sess.Get("AdminData")

	// Get query parameters
	searchKey := c.Query("searchkey", "")
	searchBy := c.Query("searchby", "transaction_id")
	status := c.Query("status", "")
	date_1st := c.Query("date_1st", "")
	date_2nd := c.Query("date_2nd", "")
	time_period := c.Query("time_period", "")

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset := (page - 1) * limit

	searchString := ""
	searchStringFull := ""

	if searchKey != "" && searchBy != "" {
		searchString = " AND " + searchBy + " ILIKE " + "'%" + searchKey + "%'"
		searchStringFull = searchStringFull + "" + searchString
	}

	if status != "" {
		searchString = " AND status ILIKE " + "'%" + status + "%'"
		searchStringFull = searchStringFull + "" + searchString
	}

	if date_1st != "" && date_2nd != "" {
		searchString = " AND createdate >= " + "'" + date_1st + "' AND createdate <= " + "'" + date_2nd + "'"
		searchStringFull = searchStringFull + "" + searchString
	}

	if len(searchStringFull) > 4 {
		searchStringFull = searchStringFull[4:]
	}

	fmt.Println("searchString => ", searchStringFull)

	transactionList := []models.Transaction_Pay{}
	database.DB.Db.Table("transaction").Order("id desc").Where(searchStringFull).Limit(limit).Offset(offset).Find(&transactionList)

	var total int64
	database.DB.Db.Table("transaction").Where(searchStringFull).Count(&total)

	// Display coin list in List box
	coinList := []models.CoinList{}
	database.DB.Db.Table("coin_list").Order("coin ASC").Where("status = ?", 1).Find(&coinList)

	//fmt.Println(total)

	Alerts := sess.Get("AlertX")
	if Alerts != "" {
		sess.Delete("AlertX")
		if err := sess.Save(); err != nil {
			panic(err)
		}
	}

	//fmt.Println(transactionList)
	return c.Render("admin/transactions", fiber.Map{
		"Title":           "Transactions",
		"Subtitle":        "Transactions",
		"AlertX":          Alerts,
		"AdminData":       adminData,
		"TransactionList": transactionList,
		"CoinList":        coinList,
		"Page":            page,
		"Limit":           limit,
		"Total":           total,
		"SearchKey":       searchKey,
		"SearchBy":        searchBy,
		"Status":          status,
		"Date_1st":        date_1st,
		"Date_2nd":        date_2nd,
		"Time_period":     time_period,
	})
}

// Function for Display Transaction Details in admin section
func AdminTransDetailsView(c *fiber.Ctx) error {

	AdminSession(c)
	sess, _ := store.Get(c)
	adminData := sess.Get("AdminData")

	Alerts := sess.Get("AlertX")
	if Alerts != "" {
		sess.Delete("AlertX")
		if err := sess.Save(); err != nil {
			panic(err)
		}
	}

	TID := c.Query("TID")

	transData := models.Transaction_Pay{}
	database.DB.Db.Table("transaction").Where("transaction_id = ?", TID).Find(&transData)

	//=============Fetch coin list ===============
	coinList := []models.CoinList{}
	database.DB.Db.Table("coin_list").Order("coin ASC").Where("status = ?", 1).Find(&coinList)

	addressUrl := models.CoinAddressUrl{}
	database.DB.Db.Table("coin_list").Select("coin_status_url").Where("coin = ?", strings.ToLower(transData.Receivedcurrency)).Find(&addressUrl)

	return c.Render("admin/trans-details", fiber.Map{
		"Title":      "Transactions Details",
		"Subtitle":   "Transactions Details",
		"TransData":  transData,
		"AdminData":  adminData,
		"CoinList":   coinList,
		"AddressUrl": addressUrl.Coin_status_url,
		"AlertX":     Alerts,
	})
}

// Function for Post Transaction Approval in admin section
func AdminTransApprovalPost(c *fiber.Ctx) error {

	////////////////////////////////////
	sess, _ := store.Get(c)

	// Get data from Form
	Alerts := "Transaction Update Successfully"
	tableID := c.FormValue("tableID")
	// convert string to Uint value
	cid, err := strconv.ParseUint(tableID, 10, 32)
	if err != nil {
		fmt.Println("Error 105")
	}
	getTableID := uint(cid)
	txhash := c.FormValue("txhash")
	status := c.FormValue("status")
	trackID := c.FormValue("trackID")

	approver_comment := c.FormValue("approver_comment")
	receivedamount := c.FormValue("receivedamount")

	// convert string to float value
	recvdmt, err := strconv.ParseFloat(receivedamount, 64)
	if err != nil {
		fmt.Println(err)
	}
	// Fetch Data from Admin Session
	adminData := sess.Get("AdminData")

	if adminData == nil {
		return c.Redirect("/admin/login")
	}
	// Convert the session data to a map
	adminMap := adminData.(map[string]interface{})
	//fmt.Println(adminMap["AdminID"])
	AdminID := adminMap["AdminID"].(uint)
	AdminName := adminMap["AdminName"].(string)

	currentTime := time.Now()
	// Format the current time as a string
	formattedTime := currentTime.Format("2006-01-02 15:04:05")

	if trackID != "" {
		database.DB.Db.Table("transaction").Save(&models.Transaction_Withdraw_Update{Id: getTableID, Response_hash: txhash, Approved_by: AdminName, Approver_comment: approver_comment, Approved_date: formattedTime, Status: status, Approver_id: AdminID})
		database.DB.Db.Table("transaction").Where("customerrefid = ?", trackID).Updates(&models.Transaction_Withdraw_Update{Approved_by: AdminName, Approver_comment: approver_comment, Approved_date: formattedTime, Status: status, Approver_id: AdminID}).Omit("id")

	} else {

		database.DB.Db.Table("transaction").Save(&models.Transaction_Update{Id: getTableID, Receivedamount: recvdmt, Response_hash: txhash, Approved_by: AdminName, Approver_comment: approver_comment, Approved_date: formattedTime, Status: status, Approver_id: AdminID})
	}

	///////////////////////==================
	sess.Set("AlertX", Alerts)
	if err := sess.Save(); err != nil {
		fmt.Println("session not store on line no 560")
	}

	return c.Redirect("/admin/transactions")
}

// Function for Display Transaction Listing in admin section
func AdminWithdrawsView(c *fiber.Ctx) error {

	AdminSession(c)
	sess, _ := store.Get(c)
	adminData := sess.Get("AdminData")

	// Get query parameters
	searchKey := c.Query("searchkey", "")
	searchBy := c.Query("searchby", "transaction_id")
	status := c.Query("status", "")
	date_1st := c.Query("date_1st", "")
	date_2nd := c.Query("date_2nd", "")
	time_period := c.Query("time_period", "")

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset := (page - 1) * limit

	searchString := ""
	searchStringFull := " AND transaction_type='WITHDRAW' "

	if searchKey != "" && searchBy != "" {
		searchString = " AND " + searchBy + " ILIKE " + "'%" + searchKey + "%'"
		searchStringFull = searchStringFull + "" + searchString
	}

	if status != "" {
		searchString = " AND status ILIKE " + "'%" + status + "%'"
		searchStringFull = searchStringFull + "" + searchString
	}

	if date_1st != "" && date_2nd != "" {
		searchString = " AND createdate >= " + "'" + date_1st + "' AND createdate <= " + "'" + date_2nd + "'"
		searchStringFull = searchStringFull + "" + searchString
	}

	if len(searchStringFull) > 4 {
		searchStringFull = searchStringFull[4:]
	}

	fmt.Println("searchString => ", searchStringFull)

	transactionList := []models.Transaction_Pay{}
	database.DB.Db.Table("transaction").Order("id desc").Where(searchStringFull).Limit(limit).Offset(offset).Find(&transactionList)

	var total int64
	database.DB.Db.Table("transaction").Where(searchStringFull).Count(&total)

	// Display coin list in List box
	coinList := []models.CoinList{}
	database.DB.Db.Table("coin_list").Order("coin ASC").Where("status = ?", 1).Find(&coinList)

	//fmt.Println(total)

	Alerts := sess.Get("AlertX")
	if Alerts != "" {
		sess.Delete("AlertX")
		if err := sess.Save(); err != nil {
			panic(err)
		}
	}

	//fmt.Println(transactionList)
	return c.Render("admin/withdraw-list", fiber.Map{
		"Title":           "Withdraw List",
		"Subtitle":        "Withdraw List",
		"AlertX":          Alerts,
		"AdminData":       adminData,
		"TransactionList": transactionList,
		"CoinList":        coinList,
		"Page":            page,
		"Limit":           limit,
		"Total":           total,
		"SearchKey":       searchKey,
		"SearchBy":        searchBy,
		"Status":          status,
		"Date_1st":        date_1st,
		"Date_2nd":        date_2nd,
		"Time_period":     time_period,
	})
}

// Function for Display Transaction Listing in admin section
func AdminRevenueView(c *fiber.Ctx) error {

	AdminSession(c)
	sess, _ := store.Get(c)
	adminData := sess.Get("AdminData")

	// Get query parameters
	searchKey := c.Query("searchkey", "")
	searchBy := c.Query("searchby", "transaction_id")
	status := c.Query("status", "")
	date_1st := c.Query("date_1st", "")
	date_2nd := c.Query("date_2nd", "")
	time_period := c.Query("time_period", "")

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset := (page - 1) * limit

	searchString := ""
	searchStringFull := " AND transaction_type='FEE' "

	if searchKey != "" && searchBy != "" {
		searchString = " AND " + searchBy + " ILIKE " + "'%" + searchKey + "%'"
		searchStringFull = searchStringFull + "" + searchString
	}

	if status != "" {
		searchString = " AND status ILIKE " + "'%" + status + "%'"
		searchStringFull = searchStringFull + "" + searchString
	}

	if date_1st != "" && date_2nd != "" {
		searchString = " AND createdate >= " + "'" + date_1st + "' AND createdate <= " + "'" + date_2nd + "'"
		searchStringFull = searchStringFull + "" + searchString
	}

	if len(searchStringFull) > 4 {
		searchStringFull = searchStringFull[4:]
	}

	fmt.Println("searchString => ", searchStringFull)

	transactionList := []models.Transaction_Pay{}
	database.DB.Db.Table("transaction").Order("id desc").Where(searchStringFull).Limit(limit).Offset(offset).Find(&transactionList)

	var total int64
	database.DB.Db.Table("transaction").Where(searchStringFull).Count(&total)

	// Display coin list in List box
	coinList := []models.CoinList{}
	database.DB.Db.Table("coin_list").Order("coin ASC").Where("status = ?", 1).Find(&coinList)

	//fmt.Println(total)

	Alerts := sess.Get("AlertX")
	if Alerts != "" {
		sess.Delete("AlertX")
		if err := sess.Save(); err != nil {
			panic(err)
		}
	}

	//fmt.Println(transactionList)
	return c.Render("admin/revenue", fiber.Map{
		"Title":           "Revenue",
		"Subtitle":        "Revenue",
		"AlertX":          Alerts,
		"AdminData":       adminData,
		"TransactionList": transactionList,
		"CoinList":        coinList,
		"Page":            page,
		"Limit":           limit,
		"Total":           total,
		"SearchKey":       searchKey,
		"SearchBy":        searchBy,
		"Status":          status,
		"Date_1st":        date_1st,
		"Date_2nd":        date_2nd,
		"Time_period":     time_period,
	})
}
