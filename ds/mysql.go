package ds

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type DataSourceMySQL struct {
	Conn *sql.DB
}

func (ds *DataSourceMySQL) Import() bool {
	return true
}

func (ds *DataSourceMySQL) InitConn() error {

	dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", "root", "root", "127.0.0.1", "3306", "ttfund", "utf8")
	conn, err := sql.Open("mysql", dbDSN)

	if err != nil {
		return err
	}

	ds.Conn = conn
	return nil
}

func (ds *DataSourceMySQL) GetData() string {

	mSql := fmt.Sprintf("select * from tt_fund limit 1000 offset %d", 0)
	rows, err := ds.Conn.Query(mSql)

	if err != nil {
		panic(err)
	}

	cols, _ := rows.Columns()

	scans := make([]interface{}, len(cols))
	vals := make([][]byte, len(cols))

	for k, _ := range vals {
		scans[k] = &vals[k]
	}
	i := 0
	result := make(map[int]map[string]string)
	for rows.Next() {

		rows.Scan(scans...)
		row := make(map[string]string)
		for k, v := range vals {
			key := cols[k]
			row[key] = string(v)
		}

		result[i] = row
		i++
	}
	rows.Close()

	fmt.Println(result)

	return "12"
}
