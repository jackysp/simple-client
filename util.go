package simple_client

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	THREADS = 1
	BATCH   = 1000000 / THREADS
)

func importData(db *sql.DB) {
	_, err := db.Exec("drop table if exists t")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("create table t (i int)")
	if err != nil {
		log.Fatal(err)
	}
	var wg sync.WaitGroup
	wg.Add(THREADS)
	start := time.Now()
	for i := 0; i < THREADS; i++ {
		go func(seq int) {
			txn, err := db.Begin()
			if err != nil {
				log.Fatal(err)
			}

			var stmt *sql.Stmt
			driverType := reflect.TypeOf(db.Driver())
			if driverType.String() == "*pq.Driver" {
				stmt, err = txn.Prepare("insert into t values ($1)")
			} else if driverType.String() == "*mysql.MySQLDriver" {
				stmt, err = txn.Prepare("insert into t values (?)")
			} else {
				log.Fatalf("unknown driver %s", driverType)
			}
			if err != nil {
				log.Fatal(err)
			}

			for j := BATCH *seq; j < BATCH*(seq+1); j++ {
				_, err = stmt.Exec(j)
				if err != nil {
					log.Fatal(err)
				}
			}

			err = stmt.Close()
			if err != nil {
				log.Fatal(err)
			}

			err = txn.Commit()
			if err != nil {
				log.Fatal(err)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	log.Println("import finish:", time.Since(start))
}

func testDML(db *sql.DB) {
	_, err := db.Exec("drop table if exists t1")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("create table t1 (i int, unique (i))")
	if err != nil {
		log.Fatal(err)
	}

	func () {
		start := time.Now()
		_, err = db.Exec("insert into t1 select * from t")
		//_, err = db.Exec("insert into t1 select * from t order by i")
		if err != nil {
			log.Fatal(err)
		}
		log.Println("insert select test finish:", time.Since(start))
	}()

	func () {
		start := time.Now()
		_, err = db.Exec("select count(*) from t1")
		if err != nil {
			log.Fatal(err)
		}
		log.Println("select test finish:", time.Since(start))
	}()

	/*
	func () {
		start := time.Now()
		_, err = db.Exec("update t1 set i = -i")
		if err != nil {
			log.Fatal(err)
		}
		log.Println("update test finish:", time.Since(start))
	}()

	func () {
		start := time.Now()
		_, err = db.Exec("delete from t1")
		if err != nil {
			log.Fatal(err)
		}
		log.Println("delete test finish:", time.Since(start))
	}()
	 */
}

func testInsert(db *sql.DB) {
	_, err := db.Exec("create table if not exists t2 (i int)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("truncate table t2")
	if err != nil {
		log.Fatal(err)
	}
	total := BATCH * THREADS
	values := make([]string, 0, total)
	for i := 0; i < total; i++ {
		values = append(values, "(" + strconv.Itoa(i) + ")")
	}
	sqlStr := fmt.Sprintf("insert into t2 values %s", strings.Join(values, ","))
	start := time.Now()
	_, err = db.Exec(sqlStr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("insert finish:", time.Since(start))
}
