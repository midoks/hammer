package ds

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type DataSourceMySQL struct {
}

func (ds *DataSourceMySQL) Import() bool {
	return true
}

func (ds *DataSourceMySQL) GetData() string {
	dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", "root", "root", "127.0.0.1", "3306", "ttfund", "utf8")

	mdb, _ := sql.Open("mysql", dbDSN)

	mSql := fmt.Sprintf("select * from tt_fund limit 1000 offset %d", 0)

	rows, err := mdb.Query(mSql)

	if err != nil {
		panic(err)
	}

	cols, _ := rows.Columns()

	scans := make([]interface{}, len(cols))
	vals := make([][]byte, len(cols))

	for k, _ := range vals {
		scans[k] = &vals[k]
	}

	for rows.Next() {

		rows.Scan(scans...)
		row := make(map[string]string)
		for k, v := range vals {
			key := cols[k]
			row[key] = string(v)
		}

		fmt.Println(row)
	}
	rows.Close()

	return "12"
}
