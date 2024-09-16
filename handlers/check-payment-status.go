package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"template/database"
	"template/models"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type StatusResponse struct {
	Hashcode       string `json:"hashcode"`
	Payment_status string `json:"payment_status"`
	Payment_id     string `json:"payment_id"`
}
type StatusRequest struct {
	Status_coin    string `json:"status_coin"`
	Status_address string `json:"status_address"`
	Status_transid string `json:"status_transid"`
	Status_coinid  string `json:"status_coinid"`
}

// Function for Post Payment form and display qr code in merchant section
func FetchPaymentStatus(c *fiber.Ctx) error {

	// Get Data from ajax
	req := new(StatusRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request",
		})
	}

	status_coin := req.Status_coin // for set condition
	//fmt.Println("status_coin==========>", status_coin)
	status_address := req.Status_address
	//fmt.Println("status_address==========>", status_address)
	status_transid := req.Status_transid
	//fmt.Println("status_transid ==========>", status_transid)
	status_coinid := req.Status_coinid
	//fmt.Println("Coin ID ==========>", status_coinid)

	// For Address
	coinAddress := models.AddressListing{}
	database.DB.Db.Table("coin_address").Select("token_details,lastupdate").Where("coin = ? AND address=? ", status_coin, status_address).Find(&coinAddress)

	fetchToken := coinAddress.Token_details
	//fetchTimestamp := coinAddress.Lastupdate

	//fmt.Println(fetchTimestamp)

	var tokenData models.OnlineTokenData
	if err := json.Unmarshal([]byte(fetchToken), &tokenData); err != nil {
		fmt.Println(err)
	}
	//fmt.Println(tokenData)

	fmt.Println("Contact Address: - ", tokenData.TokenId) // fetch from Address table json
	fmt.Println("Address :- ", status_address)

	receivedFrom := ""
	receivedTo := ""
	receivedHash := ""
	receivedFinalResult := ""
	fetchTimestamp := ""
	responseTimestamp := ""
	receivedAmountNew := 0.00
	var body []byte

	//Start status By Crypto ID

	if status_coinid == "17" { // For USDT Tether TRX/TRC20

		// URL of the external site to fetch JSON from
		url := "https://apilist.tronscanapi.com/api/token_trc20/transfers-with-status?limit=1&start=0&trc20Id=" + tokenData.TokenId + "&address=" + status_address
		//fmt.Println("Path :- ", path)

		//////////////////////////////////////
		resp, err := http.Get(url)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Failed to fetch data")
		}
		defer resp.Body.Close()

		// Reading the response body
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Failed to read response body")
		}
		// Initialize the Response struct
		var responseD models.OnlineResponse

		// Unmarshal the byte data into the struct
		err = json.Unmarshal(body, &responseD)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Failed to parse JSON")
		}
		//fmt.Println(" Data :", responseD)
		receivedAmount := responseD.Data[0].Amount

		// convert string to float value
		receivedAmt, err := strconv.ParseFloat(receivedAmount, 64)
		if err != nil {
			fmt.Println(" Error convert string to float value :")
		}
		// Convert the integer to a float64 with 6 decimal places
		AmountInFloat := float64(receivedAmt) / 1000000

		// Format the float to 6 decimal places as a string
		formattedResult := strconv.FormatFloat(AmountInFloat, 'f', 6, 64)

		// convert string to float value
		receivedAmountNew, err = strconv.ParseFloat(formattedResult, 64)
		if err != nil {
			fmt.Println(" Error convert string to float value :")
		}

		//fmt.Println("receivedAmountNew", receivedAmountNew)
		receivedFrom = responseD.Data[0].From
		receivedTo = responseD.Data[0].To
		receivedHash = responseD.Data[0].Hash
		receivedBlockTimestamp := responseD.Data[0].BlockTimestamp
		//receivedDecimals := responseD.Data[0].Decimals
		receivedFinalResult = responseD.Data[0].FinalResult
		//receivedEventType := responseD.Data[0].EventType

		//fetchTimestamp = "2024-08-23 16:00:09"
		//fetchTimestamp = coinAddress.LastUpdate.Format("2006-01-02 15:04:05")
		fetchTimestamp = coinAddress.Lastupdate.Format("2006-01-02 15:04:05")
		responseTimestamp = time.Unix(receivedBlockTimestamp/1000, 0).Format("2006-01-02 15:04:05")
		//fmt.Println("fetchTimestamp => ", fetchTimestamp)
		//fmt.Println("Response Time => ", responseTimestamp)

	} else {

		fmt.Println("Crypto Not Supported ==> ", status_coin)
	}
	//End $$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$//$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$//$$$$$$$$$$$$$$$$$$$$$

	// Calculate the difference
	layout := "2006-01-02 15:04:05" // Define the layout for parsing
	// Calculate the difference
	// Parse the dateTime strings to time.Time objects
	dateTime1, err := time.Parse(layout, fetchTimestamp)
	if err != nil {
		return c.Status(400).SendString("Error parsing dateTime1: " + err.Error())
	}

	dateTime2, err := time.Parse(layout, responseTimestamp)
	if err != nil {
		return c.Status(400).SendString("Error parsing dateTime2: " + err.Error())
	}

	duration := dateTime2.Sub(dateTime1)

	// Get the difference in seconds
	seconds := int(duration.Seconds())

	if seconds > 5 {
		fmt.Println("Success Transaction")
		database.DB.Db.Table("transaction").Where("transaction_id = ?", status_transid).Updates(&models.TransactionUpdateOnline{Receivedamount: receivedAmountNew, Receivedcurrency: status_coin, Response_hash: receivedHash, Response_from: receivedFrom, Response_to: receivedTo, Status: receivedFinalResult, Response_timestamp: dateTime2, Response_json: string(body)}).Omit("id")

	} else {

		fmt.Println("Failed Transaction")
		receivedHash = ""
		receivedFinalResult = "Response Waiting"
	}

	///////////////////////////////////////

	response := StatusResponse{
		Hashcode:       receivedHash,
		Payment_status: receivedFinalResult,
		Payment_id:     status_transid,
	}

	return c.JSON(response)
}
