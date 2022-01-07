package database

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/mburaksoran/models"
	"github.com/mburaksoran/shared"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", shared.Config.POSTGRESURL)
	if err != nil {
		log.Fatal(err)
	}

}

func GetOrdersWithRestID(id int) ([]*models.OrderFromDB, error) {
	rows, err := db.Query("SELECT * FROM orders WHERE rest_id =" + strconv.Itoa(id))
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No Records Found")
			return nil, err
		}
	}
	defer rows.Close()
	var orders []*models.OrderFromDB
	for rows.Next() {
		prd := &models.OrderFromDB{}
		err := rows.Scan(&prd.OrderID, &prd.OrderTime, &prd.OrderElemCount, &prd.OrderValue, &prd.RestID, &prd.UserID, &prd.ItemID)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		orders = append(orders, prd)

	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}
	return orders, nil
}
