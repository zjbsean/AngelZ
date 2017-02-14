package testlogic

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func RowBytesTest() {
	db, err := DBInit("root", "123456", "127.0.0.1", "3306", "go_test")
	if err != nil {
		return
	}
	defer db.Close()

	rows, err := db.Query("select * from squareNum")
	if err != nil {
		fmt.Println("DB Query Fail : Err=", err)
		return
	}

	columns, err := rows.Columns()
	if err != nil {
		fmt.Println("Get Columns Fail : Err=", err)
		return
	}
	fmt.Println("Columns Count : ", len(columns))
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(columns))

	for i, _ := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			fmt.Println("Rows Scan Fail : Err=", err)
			break
		}

		var value string
		for i, col := range values {
			if col == nil {
				fmt.Println("Value == nill")
			} else {
				value = string(col)
			}
			fmt.Println(columns[i], ": ", value)
		}
		fmt.Println("-----------------------------------")
	}
	if err = rows.Err(); err != nil {
		fmt.Println("Rows Err , Err=", err)
	}
}

func RowTest() {
	db, err := DBInit("root", "123456", "127.0.0.1", "3306", "go_test")
	if err != nil {
		return
	}
	defer db.Close()

	rows, err := db.Query("select * from squareNum")
	if err != nil {
		fmt.Println("DB Query Fail : Err=", err)
		return
	}

	columns, err := rows.Columns()
	if err != nil {
		fmt.Println("Get Columns Fail : Err=", err)
		return
	}
	fmt.Println("Columns Count : ", len(columns))

	for rows.Next() {
		var (
			number       int32
			squareNumber int32
			DT           time.Time
		)
		//var dtBytes []uint8
		err = rows.Scan(&number, &squareNumber, &DT)
		if err != nil {
			fmt.Println("Rows Scan Fail : Err=", err)
			break
		}
		fmt.Println(columns[0], ": ", number)
		fmt.Println(columns[1], ": ", squareNumber)

		fmt.Println(columns[2], ": ", DT)

		fmt.Println("-----------------------------------")
	}
	if err = rows.Err(); err != nil {
		fmt.Println("Rows Err , Err=", err)
	}
}
