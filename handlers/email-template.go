package handlers

import (
	"fmt"
	"strconv"
	"template/database"
	"template/models"

	"github.com/gofiber/fiber/v2"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var tableNameT = "email_template"
var pageNameT = "admin/email-template"
var pageTitleT = "Email Template"
var listSizeT = "10"
var listOrderByT = "status ASC,template_code ASC"

// For Add Edit Delete Email Template from admin section for send email with related template

// function for Display Email Template List
func GetTemplateList(c *fiber.Ctx) error {

	AdminSession(c)
	// Session Check
	sess, _ := store.Get(c)
	adminData := sess.Get("AdminData")

	Alerts := sess.Get("AlertX")
	if Alerts != "" {
		sess.Delete("AlertX")
		if err := sess.Save(); err != nil {
			panic(err)
		}
	}

	// Get query parameters for page and limit
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid page number")
	}
	limit, err := strconv.Atoi(c.Query("limit", listSizeT))
	if err != nil || limit < 1 {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid limit number")
	}

	// Calculate offset
	offset := (page - 1) * limit

	dataList := []models.EmailTemplate{}

	var total int64
	database.DB.Db.Table(tableNameT).Order(listOrderByT).Limit(limit).Offset(offset).Find(&dataList).Count(&total)

	//fmt.Println(total)
	// Prepare pagination data
	nextPage := page + 1
	prevPage := page - 1
	if page == 1 {
		prevPage = 0
	}

	return c.Render(pageNameT, fiber.Map{
		"Title":     pageTitleT,
		"Subtitle":  pageTitleT,
		"Action":    "List",
		"AlertX":    Alerts,
		"DataList":  dataList,
		"AdminData": adminData,
		"NextPage":  nextPage,
		"PrevPage":  prevPage,
		"Limit":     limit,
		"Count":     total,
	})
}

// function for Display Email Template Form
func AddTemplateView(c *fiber.Ctx) error {

	AdminSession(c)
	// Session Check
	sess, _ := store.Get(c)
	adminData := sess.Get("AdminData")

	return c.Render(pageNameT, fiber.Map{
		"Title":     pageTitleT,
		"Subtitle":  pageTitleT,
		"Action":    "Add",
		"AdminData": adminData,
	})
}

// function for Post Add / Edit Email Template Form
func TemplatePost(c *fiber.Ctx) error {

	AdminSession(c)
	// Parses the request body
	template_code := c.FormValue("template_code")
	template_subject := c.FormValue("template_subject")
	template_desc := c.FormValue("template_desc")
	status1, err := strconv.ParseInt(c.FormValue("status"), 10, 32)
	if err != nil {
		fmt.Println("Error 104")
	}
	status := int(status1)

	tableID := c.FormValue("tableID")
	cid, err := strconv.ParseUint(tableID, 10, 32)
	if err != nil {
		fmt.Println("Error 105")
	}
	getTableID := uint(cid)

	//////////
	// if GET ID than work update else insert
	// for Full path use- filePath & only file name use file.Filename
	result := database.DB.Db.Table(tableNameT).Save(&models.EmailTemplate{Template_id: getTableID, Template_code: template_code, Template_subject: template_subject, Template_desc: template_desc, Status: status})

	//fmt.Println(loginList.Status)
	Alerts := pageTitleT + " Processed successfully"
	if result.Error != nil {
		//fmt.Println("ERROR in QUERY")
		Alerts = pageTitleT + " Not Updated"
	}

	// check session
	sess, _ := store.Get(c)
	sess.Set("AlertX", Alerts) // Set a session key
	if err := sess.Save(); err != nil {
		return err
	}

	return c.Redirect("/" + pageNameT)
}

// function for Display Email Template Form for Edit
func EditTemplate(c *fiber.Ctx) error {

	AdminSession(c)
	tableID := c.Params("TID")

	// Session Check
	sess, _ := store.Get(c)
	adminData := sess.Get("AdminData")

	data := models.EmailTemplate{}
	database.DB.Db.Table(tableNameT).Where("template_id = ?", tableID).Find(&data)

	return c.Render(pageNameT, fiber.Map{
		"Title":     pageTitleT,
		"Subtitle":  pageTitleT,
		"Action":    "Edit",
		"AdminData": adminData,
		"EditData":  data,
	})
}

// function for Delete Email Template
func DeleteTemplate(c *fiber.Ctx) error {
	AdminSession(c)
	id := c.Params("TID")

	// Convert string to uint
	tableID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid number format")
	}
	// Convert uint64 to uint
	getTableID := uint(tableID)

	status := 2
	// fetch data from db
	result := database.DB.Db.Table(tableNameT).Save(&models.TemplateDeleted{Template_id: getTableID, Status: status})

	Alerts := pageTitleT + " Deleted successfully"
	if result.Error != nil {
		//fmt.Println("ERROR in QUERY")
		Alerts = pageTitleT + " Not Deleted"
	}

	// check session
	sess, _ := store.Get(c)
	sess.Set("AlertX", Alerts) // Set a session key
	if err := sess.Save(); err != nil {
		return err
	}

	return c.Redirect("/" + pageNameT)
}
