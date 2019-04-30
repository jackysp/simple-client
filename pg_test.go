package simple_client

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
)

func TestPG(t *testing.T) {
	connStr := "postgres://yusp@localhost/test?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	importData(db)
	testDML(db)
	testInsert(db)
}
