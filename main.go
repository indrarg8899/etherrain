package main

import (
	"net/http"
	"database/sql"
	"fmt"
	"time"
	"encoding/json"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type CreateOrder struct {
	Email    string `json:"email"`
	SKU      string `json:"sku"`
	Quantity int    `json:"quantity"`
}

type Order struct {
	ID        int
	OrderNo   string
	Email     string
	SKU       string
	Quantity  int
	CreatedAt time.Time
}

func createMysql() (*sql.DB, error) {
	dataSourceName := fmt.Sprintf(
		"%s:%s@tcp(%s:3306)/hello_go_db?parseTime=true",
		"root",
		"",
		"localhost")
	sqlDB, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}
	return sqlDB, nil
}

const (
	insertQuery = `
		INSERT INTO orders(
			order_number,
			email,
			sku,
			quantity,
			created_at) VALUES (?, ?, ?, ?, ?)
		`
)

func insertRow(sqlDB *sql.DB, order Order) error {
	_, err := sqlDB.Exec(insertQuery,
		order.OrderNo,
		order.Email,
		order.SKU,
		order.Quantity,
		order.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func main()  {
	sqlDB, err := createMysql()
	if err != nil {
		fmt.Println(err)
	}
	_, err = sqlDB.Exec(insertQuery,
		"012",
		"sarascahya@live.com",
		"BALI-I",
		12,
		time.Now())

	if err != nil {
		fmt.Println(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/order/{id}", func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		orderID := vars["id"]
		row := sqlDB.QueryRow(
			"SELECT * from orders where id=?", 
			orderID)
		var order Order
		err = row.Scan(
				&order.ID,
				&order.OrderNo,
				&order.Email,
				&order.SKU,
				&order.Quantity,
				&order.CreatedAt)
		if err != nil {
				w.Write([]byte(err.Error()))
				return
			}
		plainInfo := fmt.Sprintf("%s => %s", order.Email, order.OrderNo)
		w.Write([]byte(plainInfo))
	})

	router.HandleFunc("/order", func(w http.ResponseWriter, req *http.Request) {
		var createOrder CreateOrder
		err := json.NewDecoder(req.Body).
			Decode(&createOrder)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		orderEntity := Order{
			OrderNo:   "RANDOM-No",
			Email:     createOrder.Email,
			SKU:       createOrder.SKU,
			Quantity:  createOrder.Quantity,
			CreatedAt: time.Now(),
		}
		err = insertRow(sqlDB, orderEntity)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		fmt.Printf("Email %s \n", createOrder.Email)
		w.Write([]byte(createOrder.Email))
	}).Methods("POST")
	http.ListenAndServe(":8080", router)
}


