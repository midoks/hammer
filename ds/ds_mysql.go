package ds

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type DataSourceMySQL struct {
	Conn     *sql.DB
	DataChan chan map[int]map[string]string
}

func (ds *DataSourceMySQL) getPage(p int) (map[int]map[string]string, error) {
	result := make(map[int]map[string]string)

	mSql := fmt.Sprintf("select * from tt_fund limit 3 offset %d", p)
	fmt.Println(mSql)
	rows, err := ds.Conn.Query(mSql)

	if err != nil {
		return result, err
	}

	cols, _ := rows.Columns()

	scans := make([]interface{}, len(cols))
	vals := make([][]byte, len(cols))

	for k, _ := range vals {
		scans[k] = &vals[k]
	}

	i := 0
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
	return result, nil
}

func (ds *DataSourceMySQL) Import() {

	i := 0
	for {

		result, err := ds.getPage(i)
		if err != nil {
			break
		}

		if len(result) == 0 {
			break
		}

		ds.DataChan <- result
		i++
	}
}

func (ds *DataSourceMySQL) Task() {
	for {
		d := <-ds.DataChan
		fmt.Println("Task", d)
	}
}

func (ds *DataSourceMySQL) Init() {
	ds.DataChan = make(chan map[int]map[string]string, 1000)
	ds.InitConn()
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

func (ds *DataSourceMySQL) GetData() (map[int]map[string]string, error) {

	result := make(map[int]map[string]string)

	mSql := fmt.Sprintf("select * from tt_fund limit 1 offset %d", 0)
	rows, err := ds.Conn.Query(mSql)

	if err != nil {
		return result, err
	}

	cols, _ := rows.Columns()

	scans := make([]interface{}, len(cols))
	vals := make([][]byte, len(cols))

	for k, _ := range vals {
		scans[k] = &vals[k]
	}

	i := 0
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

	return result, nil
}
