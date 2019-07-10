package simple_client

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestMysql(t *testing.T) {
	connStr := "root@tcp(localhost)/test?maxAllowedPacket=512000000"
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	importData(db)
	testDML(db)
	//testInsert(db)
}

