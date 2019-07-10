package simple_client

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestTiDB(t *testing.T) {
	connStr := "root@tcp(172.16.5.34:8001)/test?maxAllowedPacket=512000000"
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	importData(db)
	//testDML(db)
	//testInsert(db)
}

