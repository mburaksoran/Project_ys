package utils

import (
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/mburaksoran/models"
)

func UploadData(r *http.Request) string {
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, _, err := r.FormFile("myFile")
	if err != nil {
		// todo proper error handling
		fmt.Println("Error Retrieving the File")

	}
	defer file.Close()

	// Create a temporary file within our temp directory that follows
	// a particular naming pattern
	tempFile, err := ioutil.TempFile("temp", "data-*.csv")

	if err != nil {
		//todo proper error handling
		fmt.Println(err)
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		//todo proper error handling
		fmt.Println(err)
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)

	//mod := UploadedDataoStruct(tempFile.Name())
	return tempFile.Name()
}

func UploadedDataToStruct(filepath string) []models.OrderFromClient {
	// for reading and turning into list of Customer struct, first needed variable created.

	var order models.OrderFromClient
	var orders []models.OrderFromClient
	// Creating a connection via csv files
	csvFile, err := os.Open(filepath)
	defer csvFile.Close()

	if err != nil {
		fmt.Println(err)
	}
	//creating a *csv.reader
	reader := csv.NewReader(csvFile)
	//reading csv files via reader and method of reader
	csvData, err := reader.ReadAll()
	if err != nil {
		// TODO proper error handling
		fmt.Println(err)
		os.Exit(1)

	}
	//in csv file, each comma represent a column, with this knowledge every row has 9 column and each column parse where they are belong.

	for i, each := range csvData {
		if i > 0 {
			orderID, _ := strconv.Atoi(each[0])
			order.OrderID = orderID
			order.SentBy = each[1]
			order.OrderTime = each[2]
			itemval, _ := strconv.ParseFloat(each[3], 64)
			order.ItemVal = itemval
			order.ClosedBy = each[4]
			order.RestName = each[5]
			order.PaidType = each[6]
			order.From = each[7]
			order.Item = each[8]
			itemcount, _ := strconv.Atoi(each[9])
			order.ItemCount = itemcount
			order.TableName = each[10]
			order.ClosedAt = each[11]

			orders = append(orders, order)
		}
		// first object of column is userID and it type actually integer but when reader read csv files it turns to string. in this line string converted to integer via strconv package.

		// checking all the elements already readed. if id is match with existed struct elements, only orderamounth and order discont will added to existed element but if there is no match
		// function will create new struct elements.

	}

	return orders
}

func DataPrepForClient(client []models.OrderFromClient, database []*models.OrderFromDB) []models.DataForClient {
	var temp models.DataForClient
	var templist []models.DataForClient
	for _, i := range client {
		temp.OrderTime = i.ClosedAt
		temp.DataComeFrom = 1
		temp.TotalPayment = i.ItemVal
		templist = append(templist, temp)
	}
	for _, i := range database {
		temp.OrderTime = i.OrderTime
		temp.DataComeFrom = 2
		temp.TotalPayment = float64(i.OrderValue)
		templist = append(templist, temp)
	}
	return templist
}

func PassToSHA256(passStr string) string {
	hash := sha256.New()
	hash.Write([]byte(passStr))
	sha256Hash := hex.EncodeToString(hash.Sum(nil))

	return sha256Hash
}

func GetRestIDFromToken(r *http.Request, jwtKey []byte) (int, string) {
	var claim models.Claims
	cookie, _ := r.Cookie("token")
	tokenStr := cookie.Value
	//token, _ :=
	jwt.ParseWithClaims(tokenStr, &claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	//fmt.Println(token)
	return claim.UserRestID, claim.Username
}
