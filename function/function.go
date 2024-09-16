package function

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"template/database"
	"template/models"
	"time"

	"github.com/go-resty/resty/v2"
)

// Function for send email
func SendEmail(template_code string, p models.EmailData) error {
	//////////////////////////////////////////

	//Define Variables
	Email := "mailers@itio.in"
	FullName := "Null"
	UserName := "Null"
	Password := "Null"
	Invoice_id := ""
	Invoice_url := ""
	Invoice_url_html := ""
	Amount := ""
	Crypto := ""
	Status := ""
	HashCode := ""
	TransID := ""

	// Set Variables
	if p.Email != "" {
		Email = p.Email
	}
	if p.FullName != "" {
		FullName = p.FullName
	}
	if p.Amount != "" {
		Amount = p.Amount
	}
	if p.UserName != "" {
		UserName = p.UserName
	}
	if p.Crypto != "" {
		Crypto = p.Crypto
	}
	if p.Status != "" {
		Status = p.Status
	}
	if p.HashCode != "" {
		HashCode = p.HashCode
	}
	if p.TransID != "" {
		TransID = p.TransID
	}
	if p.Password != "" {
		Password = p.Password
	}
	if p.Invoice_id != "" {
		Invoice_id = p.Invoice_id
	}
	if p.Invoice_url != "" {
		Invoice_url = p.Invoice_url
		Invoice_url_html = "<a href='" + Invoice_url + "'  target='_blank' title='Pay Now'>Pay Now</a>"
	}

	/////////////////////////////////////////////
	// Get Email Template from  table - email_template
	emailTemplate := models.EmailTemplate{}
	result := database.DB.Db.Table("email_template").Where("template_code = ?", template_code).Find(&emailTemplate)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	// smtp - Details Get from env file
	var SMTPusername = os.Getenv("SMTPusername")
	var SMTPpassword = os.Getenv("SMTPpassword")
	var fromName = os.Getenv("HostName")
	var fromEmail = os.Getenv("FromEmail")
	host := os.Getenv("SMTPhost")
	port := os.Getenv("SMTPport")
	address := host + ":" + port

	// Replace content with email template format
	subject := strings.Replace(emailTemplate.Template_subject, "[sitename]", os.Getenv("SITENAME"), 1)
	subject = strings.Replace(subject, "[fullname]", FullName, 1)
	subject = strings.Replace(subject, "[username]", UserName, 1)
	subject = strings.Replace(subject, "[status]", Status, 1)

	EmailBody := strings.Replace(emailTemplate.Template_desc, "[sitename]", os.Getenv("SITENAME"), 1)
	EmailBody = strings.Replace(EmailBody, "[fullname]", FullName, 1)
	EmailBody = strings.Replace(EmailBody, "[username]", UserName, 1)
	EmailBody = strings.Replace(EmailBody, "[password]", Password, 1)
	EmailBody = strings.Replace(EmailBody, "[invoiceid]", Invoice_id, 1)
	EmailBody = strings.Replace(EmailBody, "[invoiceurl]", Invoice_url_html, 1)
	EmailBody = strings.Replace(EmailBody, "[amount]", Amount, 1)
	EmailBody = strings.Replace(EmailBody, "[crypto]", Crypto, 1)
	EmailBody = strings.Replace(EmailBody, "[sitename1]", os.Getenv("SITENAME"), 1)
	EmailBody = strings.Replace(EmailBody, "[sitename2]", os.Getenv("SITENAME2"), 1)
	EmailBody = strings.Replace(EmailBody, "[status]", Status, 1)
	EmailBody = strings.Replace(EmailBody, "[hashCode]", HashCode, 1)
	EmailBody = strings.Replace(EmailBody, "[transid]", TransID, 1)

	to := []string{Email}
	// Set up authentication information.
	auth := smtp.PlainAuth("", SMTPusername, SMTPpassword, host)
	// create MessageBody
	msg := []byte(
		"From: " + fromName + ": <" + fromEmail + ">\r\n" +
			"To: " + Email + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME: MIME-version: 1.0\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
			"\r\n" +
			EmailBody)
	err := smtp.SendMail(address, auth, fromEmail, to, msg)
	if err != nil {
		return err
	}

	return nil
}

// Function for generate random password
func PasswordGenerator(passwordLength int) string {
	// Character sets for generating passwords
	lowerCase := "abcdefghijklmnopqrstuvwxyz" // lowercase
	upperCase := "ABCDEFGHIJKLMNOPQRSTUVWXYZ" // uppercase
	numbers := "0123456789"                   // numbers
	specialChar := "!@#$%^&*()_-+={}[/?]"     // special characters

	// Variable for storing password
	password := ""

	// Initialize the random number generator
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	// Generate password character by character
	for n := 0; n < passwordLength; n++ {
		// Generate a random number to choose a character set
		randNum := rng.Intn(4)

		switch randNum {
		case 0:
			randCharNum := rng.Intn(len(lowerCase))
			password += string(lowerCase[randCharNum])
		case 1:
			randCharNum := rng.Intn(len(upperCase))
			password += string(upperCase[randCharNum])
		case 2:
			randCharNum := rng.Intn(len(numbers))
			password += string(numbers[randCharNum])
		case 3:
			randCharNum := rng.Intn(len(specialChar))
			password += string(specialChar[randCharNum])
		}
	}

	return password
}

// GenerateRandomID generates a random ID of the specified length
func GenerateRandomID(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

type CurrencyResponse map[string]map[string]float64

// convertCurrencyToCrypto fetches the current crypto exchange rate and converts the given amount of fiat currency to cryptocurrency
func ConvertCurrencyToCrypto(amount, fromCurrency, toCrypto string) (float64, error) {
	client := resty.New()
	apiURL := "https://api.coingecko.com/api/v3/simple/price?ids=" + toCrypto + "&vs_currencies=" + fromCurrency

	//fmt.Print(apiURL)

	var apiResponse CurrencyResponse
	_, err := client.R().SetResult(&apiResponse).Get(apiURL)
	if err != nil {
		return 0, err
	}

	usdAmount, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return 0, err
	}

	rate := apiResponse[toCrypto][fromCurrency]
	cryptoAmount := usdAmount / rate

	return cryptoAmount, nil
}

type APIKeyResponse struct {
	Client_id int
}

// GetMIDByApikey Get Merchant ID from API Key
func GetMIDByApikey(apikey string) (uint, string) {

	fmt.Println("apikey - > ", apikey)

	aPIKeyResponse := APIKeyResponse{}
	database.DB.Db.Table("client_api").Select("client_id").Where("apikey = ? AND status = ?", apikey, 1).Find(&aPIKeyResponse)
	//fmt.Println("aPIKeyResponse - > ", aPIKeyResponse.Client_id)
	uid := uint(aPIKeyResponse.Client_id)
	errorx := ""
	if aPIKeyResponse.Client_id > 0 {
		return uid, errorx
	} else {
		errorx = "Wrong API Key"
		return uid, errorx
	}

}

type NameByMID struct {
	Full_name string
}

// Get Merchant Name ID  By Merchant ID
func GetNameByMID(MID uint) string {
	nameByMID := NameByMID{}
	database.DB.Db.Table("client_master").Select("full_name").Where("client_id = ? AND status = ?", MID, 1).Find(&nameByMID)
	full_name := string(nameByMID.Full_name)
	return full_name
}

type EmailByMID struct {
	Username string
}

// Get Merchant Name ID  By Merchant ID
func GetEmailByMID(MID uint) string {
	emailByMID := EmailByMID{}
	database.DB.Db.Table("client_master").Select("username").Where("client_id = ? AND status = ?", MID, 1).Find(&emailByMID)
	username := string(emailByMID.Username)
	return username
}
