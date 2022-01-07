package utils

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/mburaksoran/models"
)

func UploadData(r *http.Request) (string, error) {
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	file, _, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		return "", err
	}
	defer file.Close()

	// Create a temporary file within our temp directory that follows
	// a particular naming pattern
	tempFile, err := ioutil.TempFile("temp", "data-*.csv")
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	// write this byte array to our temporary file
	_, err = tempFile.Write(fileBytes)

	//mod := UploadedDataoStruct(tempFile.Name())
	return tempFile.Name(), err
}

func UploadedDataToStruct(filepath string) ([]models.OrderFromClient, error) {
	// for reading and turning into list of orders struct, first needed variable created.
	var order models.OrderFromClient
	var orders []models.OrderFromClient
	// Creating a connection via csv files
	csvFile, err := os.Open(filepath)
	defer csvFile.Close()

	if err != nil {
		return nil, err
	}
	//creating a *csv.reader
	reader := csv.NewReader(csvFile)
	//reading csv files via reader and method of reader
	csvData, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	//in csv file, each comma represent a column, with this knowledge every row has 12 column and each column parse where they are belong.
	for i, each := range csvData {
		if i > 0 {
			// first object of column is userID and it type actually integer but when reader read csv files it turns to string.
			//in this line string converted to integer via strconv package.
			orderID, err := strconv.Atoi(each[0])
			if err != nil {
				return nil, err
			}
			order.OrderID = orderID
			order.SentBy = each[1]
			order.OrderTime = each[2]
			itemval, err := strconv.ParseFloat(each[3], 64)
			if err != nil {
				return nil, err
			}
			order.ItemVal = itemval
			order.ClosedBy = each[4]
			order.RestName = each[5]
			order.PaidType = each[6]
			order.From = each[7]
			order.Item = each[8]
			itemcount, err := strconv.Atoi(each[9])
			if err != nil {
				return nil, err
			}
			order.ItemCount = itemcount
			order.TableName = each[10]
			order.ClosedAt = each[11]

			orders = append(orders, order)
		}

	}

	return orders, err
}

// in this project two different data source will use. one of them came from client as a csv file , the other came from client's databases data. both of them are different each other.
// to analyzing all the data, it needs to be combine and same struct. for this purpose DataPrepForClient method created. it reads both data's struct element and combine the same spesification as a one struct.
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
