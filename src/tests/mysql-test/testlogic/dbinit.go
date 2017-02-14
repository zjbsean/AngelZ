package testlogic

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func DBInit(user, passward, ip, port, dbName string) (*sql.DB, error) {
	dbParams := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=CST", user, passward, ip, port, dbName)
	db, err := sql.Open("mysql", dbParams)

	if err != nil {
		fmt.Println("Open DB Fail : ErrorMsg=", err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		defer db.Close()
		fmt.Println("Ping DB Fail : ErrorMsg=", err)
		return nil, err
	}
	fmt.Printf("Connect DB Succ ! IP=%s, Port=%s, DBName=%s\n", ip, port, dbName)
	return db, nil
}
