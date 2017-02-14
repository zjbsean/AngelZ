package testlogic

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func PreparedTest() {
	db, err := DBInit("root", "123456", "127.0.0.1", "3306", "go_test")
	if err != nil {
		return
	}
	defer db.Close()

	stmtIns, err := db.Prepare("INSERT INTO squareNum VALUES( ?, ?, ?)")
	if err != nil {
		fmt.Println("DB Prepare Insert Fail !")
		return
	}
	defer stmtIns.Close()

	stmtOut, err := db.Prepare("SELECT squareNumber FROM squarenum WHERE number = ?")
	if err != nil {
		fmt.Println("DB Prepare Select Fail !")
		return
	}
	defer stmtOut.Close()

	for i := 0; i < 25; i++ {
		_, err := stmtIns.Exec(i, (i * i), time.Now())
		if err != nil {
			fmt.Printf("DB Insert Exec Fail !, %d, %d, err=%s\n", i, i*i, err)
			break
		}
	}

	var squareNum int
	err = stmtOut.QueryRow(13).Scan(&squareNum)
	if err != nil {
		fmt.Println("DB Select 13 Fail ! err=", err)
	} else {
		fmt.Println("DB Select 13 Succ Num=", squareNum)
	}

}
