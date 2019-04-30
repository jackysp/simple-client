package simple_client

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestTiDB(t *testing.T) {
	connStr := "root@tcp(127.0.0.1:4000)/test"
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("drop table if exists t")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec("create table t (i varchar(20))")
	if err != nil {
		t.Fatal(err)
	}
	st, err := db.Prepare("insert into t values (sleep(?))")
	if err != nil {
		t.Fatal(err)
	}
	st.Exec(1)
}
