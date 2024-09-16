package handlers

import (
	"encoding/gob"
	"fmt"
	"strconv"
	"template/database"
	"template/function"
	"template/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// functions for login/ Logout / check session/ Dashboard listing/ login History/ Member / support ticket

func init() {
	gob.Register(map[string]interface{}{})
	gob.Register(string(""))
	gob.Register(int(0))
	gob.Register(float64(0))
}

// function for display admin login form
func AdminLoginView(c *fiber.Ctx) error {

	sess, err := store.Get(c)
	if err != nil {
		return err
	}
	Alerts := sess.Get("AlertX")
	if Alerts != "" {
		sess.Delete("AlertX")
		if err := sess.Save(); err != nil {
			panic(err)
		}
	}

	return c.Render("admin/login", fiber.Map{
		"Title":  "Web Admin - Sign in",
		"AlertX": Alerts,
	})
}

// function for post admin login form

func AdminLoginPost(c *fiber.Ctx) error {
	// Parses the request body
	getAdminUserName := c.FormValue("admin-username")
	getAdminPassword := c.FormValue("admin-password")

	//fmt.Println(getUserName, getPassword)
	Alerts := ""
	loginList := models.AdminLoginList{}
	result := database.DB.Db.Table("admin_master").Where("username = ?", getAdminUserName).Find(&loginList)

	//fmt.Println(loginList.Status)

	if result.Error != nil {
		//fmt.Println("ERROR in QUERY")
		Alerts = "ERROR in QUERY"
	}

	if result.RowsAffected == 1 {

		if loginList.Status != 1 {
			//fmt.Println("Account Not Activate / Deleted")
		} else if loginList.Password != "" {
			//fmt.Println(loginList.Password)
			err := bcrypt.CompareHashAndPassword([]byte(loginList.Password), []byte(getAdminPassword))
			if err == nil {
				loginIp := c.Context().RemoteIP().String()
				qry := models.LoginHistory{Client_id: loginList.Admin_id, Login_ip: loginIp, Login_type: 2}
				result := database.DB.Db.Table("login_history").Select("client_id", "login_ip", "login_type").Create(&qry)
				if result.Error != nil {
					fmt.Println(result.Error)
				}

				// Set key/value
				sess, err := store.Get(c)
				if err != nil {
					return err
				}
				sess.Set("AdminData", map[string]interface{}{
					"AdminName":  loginList.Full_name,
					"AdminEmail": getAdminUserName,
					"AdminID":    loginList.Admin_id,
					"AdminRole":  loginList.Role,
				})
				if err := sess.Save(); err != nil {
					return err
				}

				return c.Redirect("/admin/")

			} else {
				//fmt.Println("Wrong Password")
				Alerts = "Wrong Password"
			}

		}

	} else {
		//fmt.Println("Account Not Found")
		Alerts = "Account Not Found"

	}

	return c.Render("admin/login", fiber.Map{
		"Title":  "Web Admin - Sign in",
		"AlertX": Alerts,
		//"Facts":    facts,
	})
}

// Function for Display Listing and wallet list with balance
func AdminIndexView(c *fiber.Ctx) error {

	AdminSession(c)
	// // Session Check
	sess, _ := store.Get(c)
	adminData := sess.Get("AdminData")

	transactionList := []models.Transaction_Pay{}
	var total int64
	database.DB.Db.Table("transaction").Where("status = ?", "SUCCESS").Order("id desc").Limit(10).Find(&transactionList).Count(&total)

	assetList := []models.CoinWithBalance{}
	database.DB.Db.Table("transaction").Select("assetid, SUM(receivedamount)  as balance").Where("status = ?", "SUCCESS").Group("assetid").Having("COUNT(assetid) > ?", 0).Order("assetid ASC").Find(&assetList)
	//fmt.Println(assetList)

	// Display coin list in List box
	coinList := []models.CoinList{}
	database.DB.Db.Table("coin_list").Order("coin ASC").Where("status = ?", 1).Find(&coinList)

	return c.Render("admin/index", fiber.Map{
		"Title":           "Dashboard",
		"Subtitle":        "Home",
		"AdminData":       adminData,
		"TransactionList": transactionList,
		"CoinBalanceList": assetList,
		"CoinList":        coinList,
		"CountList":       total,
	})
}

// Admin Logout Function for session destroy
func AdminLogOut(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		panic(err)
	}
	sess.Delete("AdminData")

	// Destroy session
	if err := sess.Destroy(); err != nil {
		panic(err)
	}
	// Clear session or cookies
	cookie := new(fiber.Cookie)
	cookie.Name = "session_id"
	cookie.Expires = time.Now().Add(-1 * time.Hour)
	c.Cookie(cookie)

	return c.Redirect("/admin/login")
}

// Check admin Session exist or not
func AdminSession(c *fiber.Ctx) error {
	// Session Check
	sess, _ := store.Get(c)
	adminData := sess.Get("AdminData")

	if adminData == nil {
		sess.Set("AlertX", "Error - Session expired") // Set a session key
		if err := sess.Save(); err != nil {
			fmt.Println("check data - ", err)
		}
		return c.Redirect("/admin/login")
	}

	return nil
}

// function for admin login history
func AdminLoginHistory(c *fiber.Ctx) error {

	AdminSession(c)
	// // Session Check
	sess, _ := store.Get(c)
	adminData := sess.Get("AdminData")
	//fmt.Println("check data - ", adminData)
	if adminData == nil {
		return c.Redirect("/admin/login")
	}
	// Convert the session data to a map
	adminMap := adminData.(map[string]interface{})
	//fmt.Println(adminMap["AdminID"])
	LoginAdminID := adminMap["AdminID"].(uint)
	// Get query parameters for page and limit
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid page number")
	}
	limit, err := strconv.Atoi(c.Query("limit", "50"))
	if err != nil || limit < 1 {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid limit number")
	}

	// Calculate offset
	offset := (page - 1) * limit

	loginHistory := []models.LoginHistory{}
	database.DB.Db.Table("login_history").Order("token_id desc").Limit(limit).Offset(offset).Where(&models.LoginHistory{Login_type: 2, Client_id: LoginAdminID}).Find(&loginHistory)

	var total int64
	database.DB.Db.Table("login_history").Where(&models.LoginHistory{Login_type: 2, Client_id: LoginAdminID}).Count(&total)

	//fmt.Println(total)
	//fmt.Println(loginHistory)
	// Prepare pagination data
	nextPage := page + 1
	prevPage := page - 1
	if page == 1 {
		prevPage = 0
	}

	//fmt.Println(loginHistory)
	return c.Render("admin/login-history", fiber.Map{
		"Title":        "Login History",
		"Subtitle":     "Login History",
		"LoginHistory": loginHistory,
		"AdminData":    adminData,
		"NextPage":     nextPage,
		"PrevPage":     prevPage,
		"Limit":        limit,
		"Count":        total,
	})
}

// Function for Member Listing in admin section
func AdminMembersView(c *fiber.Ctx) error {

	AdminSession(c)
	sess, _ := store.Get(c)
	adminData := sess.Get("AdminData")

	// Get query parameters
	search := c.Query("search", "")
	sortBy := c.Query("sort_by", "transaction_id")
	status := c.Query("status", "")
	order := c.Query("order", "desc")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset := (page - 1) * limit

	searchString := ""
	searchStringFull := ""

	if search != "" && sortBy != "" {
		searchString = " AND " + sortBy + " ILIKE " + "'%" + search + "%'"
		searchStringFull = searchStringFull + "" + searchString
	}

	if status != "" {
		searchString = " AND status = " + status
		searchStringFull = searchStringFull + "" + searchString
	}

	if len(searchStringFull) > 4 {
		searchStringFull = searchStringFull[4:]
	}

	memberList := []models.MemberDetails{}
	//database.DB.Db.Table("client_master").Order("client_id desc").Limit(limit).Offset(offset).Find(&memberList)
	database.DB.Db.Table("client_master as a ").Where(searchStringFull).Select("a.client_id, a.Username, a.Full_name,a.status,a.timestamp, b.title, b.gender, b.country_code, b.mobile, b.address_line1, b.address_line2").Joins("left join client_details as b on b.client_id = a.client_id").Order("a.client_id ASC").Limit(limit).Offset(offset).Find(&memberList)

	var total int64
	database.DB.Db.Table("client_master").Where(searchStringFull).Count(&total)
	//fmt.Println(memberList)

	feeList := []models.FeesDetails{}
	database.DB.Db.Table("client_fees").Order("client_id ASC").Find(&feeList)

	Alerts := sess.Get("AlertX")
	if Alerts != "" {
		sess.Delete("AlertX")
		if err := sess.Save(); err != nil {
			panic(err)
		}
	}

	//fmt.Println(transactionList)
	return c.Render("admin/members", fiber.Map{
		"Title":      "Registered Members",
		"Subtitle":   "Registered Members",
		"AlertX":     Alerts,
		"AdminData":  adminData,
		"MemberList": memberList,
		"FeeList":    feeList,
		"Page":       page,
		"Limit":      limit,
		"Total":      total,
		"Search":     search,
		"SortBy":     sortBy,
		"Order":      order,
		"status":     status,
	})
}

// Function for Member Details in admin section
func AdminMembersDetailsView(c *fiber.Ctx) error {

	AdminSession(c)
	sess, _ := store.Get(c)
	adminData := sess.Get("AdminData")
	MID := c.Params("MID")

	memberList := models.MemberDetails{}
	database.DB.Db.Table("client_master as a ").Select("a.client_id, a.Username, a.Full_name,a.status,a.timestamp, b.title, b.gender, b.country_code, b.mobile, b.address_line1, b.address_line2").Joins("join client_details as b on b.client_id = a.client_id AND b.client_id=?", MID).Find(&memberList)

	feeList := models.FeesDetails{}
	database.DB.Db.Table("client_fees").Where("client_id = ? ", MID).Find(&feeList)

	cryptoWalletList := []models.CryptoWalletList{}
	database.DB.Db.Table("client_wallet").Order("status ASC,crypto_code ASC").Where("client_id = ? ", MID).Find(&cryptoWalletList)

	clientAPI := []models.ClientAPI{}
	database.DB.Db.Table("client_api").Order("id desc").Where("status = ? AND Client_id=?", 1, MID).Find(&clientAPI)

	assetList := []models.CoinWithBalance{}
	database.DB.Db.Table("transaction").Select("assetid, SUM(receivedamount)  as balance").Where("client_id = ? AND status = ?", MID, "SUCCESS").Group("assetid").Having("COUNT(assetid) > ?", 0).Order("assetid ASC").Find(&assetList)

	// Display coin list in List box
	coinList := []models.CoinList{}
	database.DB.Db.Table("coin_list").Order("coin ASC").Where("status = ?", 1).Find(&coinList)

	// fetch query for transaction list
	transactionList := []models.Transaction_Pay{}
	database.DB.Db.Table("transaction").Order("id desc").Where("status = ? AND client_id = ?", "SUCCESS", MID).Limit(50).Find(&transactionList)
	// fetch query for transaction list
	withdrawList := []models.Transaction_Pay{}
	database.DB.Db.Table("transaction").Order("id desc").Where("transaction_type = ? AND client_id = ?", "WITHDRAW", MID).Limit(50).Find(&withdrawList)
	Alerts := sess.Get("AlertX")
	if Alerts != "" {
		sess.Delete("AlertX")
		if err := sess.Save(); err != nil {
			panic(err)
		}
	}

	fmt.Println(memberList)
	fmt.Println(feeList)

	return c.Render("admin/members-details", fiber.Map{
		"Title":            "Members Details",
		"Subtitle":         "Members Details",
		"AlertX":           Alerts,
		"AdminData":        adminData,
		"MemberList":       memberList,
		"FeeList":          feeList,
		"CryptoWalletList": cryptoWalletList,
		"ClientAPI":        clientAPI,
		"CoinBalanceList":  assetList,
		"CoinList":         coinList,
		"TransactionList":  transactionList,
		"WithdrawList":     withdrawList,
	})
}

// Function for Member Active / De Active
func AdminMemberStatus(c *fiber.Ctx) error {

	AdminSession(c)
	sess, _ := store.Get(c)

	// Parses the from Query

	status1, err := strconv.ParseInt(c.Query("key"), 10, 32)
	if err != nil {
		fmt.Println("Error 104")
	}
	status := int(status1)

	cid, err := strconv.ParseUint(c.Query("mid"), 10, 32)
	if err != nil {
		fmt.Println("Error 105")
	}
	mid := uint(cid)

	result := database.DB.Db.Table("client_master").Save(&models.MemberStatusUpdate{Client_id: mid, Status: status})

	Alerts := "Member Status Updated"
	if result.Error != nil {
		//fmt.Println("ERROR in QUERY")
		Alerts = "Member Status Not Updated"
	}

	sess.Set("AlertX", Alerts)
	if err := sess.Save(); err != nil {
		panic(err)
	}

	return c.Redirect("/admin/members")
}

// Function for Support Ticket Listing - Admin
func SupportTicketListing(c *fiber.Ctx) error {

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

	ticketList := []models.Support_Ticket{}
	database.DB.Db.Table("support-ticket").Order("ticket_id desc").Limit(limit).Offset(offset).Find(&ticketList)

	var total int64
	database.DB.Db.Table("support-ticket").Count(&total)

	//fmt.Println(ticketList)

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

	//fmt.Println(transactionList)
	return c.Render("admin/support-ticket", fiber.Map{
		"Title":      "Support Ticket",
		"Subtitle":   "Support Ticket",
		"AlertX":     Alerts,
		"AdminData":  adminData,
		"TicketList": ticketList,
		"NextPage":   nextPage,
		"PrevPage":   prevPage,
		"Limit":      limit,
		"Count":      total,
	})
}

// function for Post Add / Edit Fees
func FeesPost(c *fiber.Ctx) error {

	AdminSession(c)
	// Parses the request body
	min_amount_fee := c.FormValue("min_amount_fee")
	amount_fee_in_percent := c.FormValue("amount_fee_in_percent")

	clientID, err := strconv.ParseUint(c.FormValue("client_id"), 10, 32)
	if err != nil {
		fmt.Println("Error 1041")
	}
	client_id := uint(clientID)
	client_idx := int(clientID)

	checkList := models.FeesDetails{}
	result := database.DB.Db.Table("client_fees").Where("client_id = ?", client_id).Find(&checkList)

	fmt.Println(checkList.Client_id)

	if result.RowsAffected == 1 {
		fmt.Print("Update")
		result = database.DB.Db.Table("client_fees").Save(&models.FeesDetails{Client_id: client_id, Min_amount_fee: min_amount_fee, Amount_fee_in_percent: amount_fee_in_percent}).Where("client_id = ?", client_id)

	} else {
		fmt.Print("Add")

		result = database.DB.Db.Table("client_fees").Save(&models.FeesDetail{Min_amount_fee: min_amount_fee, Amount_fee_in_percent: amount_fee_in_percent, Client_id: client_idx})
	}
	//////////

	Alerts := " Fees Processed successfully"
	if result.Error != nil {
		fmt.Println("ERROR in QUERY", result.Error)
		Alerts = "Fees Not Updated"
	}

	// check session
	sess, _ := store.Get(c)
	sess.Set("AlertX", Alerts) // Set a session key
	if err := sess.Save(); err != nil {
		return err
	}

	return c.Redirect("/admin/members")
}

// function for display admin Reset Password Form
func AdminResetPasswordView(c *fiber.Ctx) error {

	// check session
	s, _ := store.Get(c)
	Alerts := s.Get("Alertx")
	if Alerts != "" {
		s.Delete("Alertx")
		if err := s.Save(); err != nil {
			panic(err)
		}
	}

	return c.Render("admin/reset-password", fiber.Map{
		"Title":  "Reset Password",
		"Alertx": Alerts,
	})
}

// function for Submit admin Reset Password Form
func AdminResetPasswordPost(c *fiber.Ctx) error {
	// Parses the request body
	getEmail := c.FormValue("email")

	Alerts := "New Generated password has been sent to your registered Email ID"
	loginList := models.AdminList{}
	result := database.DB.Db.Table("admin_master").Where("username = ?", getEmail).Find(&loginList)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	receivedId := loginList.Admin_id
	//fmt.Println("VVV ", receivedId)

	if receivedId > 0 {
		var password = function.PasswordGenerator(8)
		var hash []byte
		// func GenerateFromPassword(password []byte, cost int) ([]byte, error)
		hash, _ = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		result := database.DB.Db.Table("admin_master").Save(&models.AdminPassword{Admin_id: loginList.Admin_id, Password: string(hash)})
		if result.Error != nil {
			//fmt.Println("ERROR in QUERY")
			Alerts = "Password Not Updated"
		}
		//  Email///
		template_code := "RESTORE-PASSWORD"

		emailData := models.EmailData{FullName: loginList.Full_name, Email: getEmail, UserName: getEmail, Password: password}

		err := function.SendEmail(template_code, emailData)
		if err != nil {
			fmt.Println("issue sending verification email")
		}

	} else {
		Alerts = "Email Not Matched with our Record. Please Check and Try with Correct"
		// check session
		sess, _ := store.Get(c)
		// check session
		sess.Set("Alertx", Alerts) // Set a session key
		if err := sess.Save(); err != nil {
			return err
		}
		return c.Redirect("/admin/reset-password")

	}

	return c.Render("admin/reset-password", fiber.Map{
		"Title":  "Reset-Password",
		"Alertx": Alerts,
	})
}
