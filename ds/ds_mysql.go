package ds

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/midoks/hammer/storage"
	"io/ioutil"
	"os"
)

type DataSourceMySQL struct {
	Conn     *sql.DB
	DataChan chan map[int]map[string]string
}

func (ds *DataSourceMySQL) getPage(p int, s int) (map[int]map[string]string, error) {
	result := make(map[int]map[string]string)

	p = s * p

	filePtr, err := os.Open("conf/test/__tmp.json")
	if err != nil {
		fmt.Println(err)
	}
	rd, err := ioutil.ReadAll(filePtr)
	if err != nil {
		fmt.Println(err)
	}
	cf := &SaveStatus{}
	_ = json.Unmarshal([]byte(rd), cf)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	fmt.Println(cf)

	mSql := fmt.Sprintf("select * from tt_fund limit %d offset %d", s, p)
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

	// i := 0
	for i := 0; i < 3; i++ {

		result, err := ds.getPage(i, 10)
		if err != nil {
			break
		}

		if len(result) == 0 {
			break
		}

		ds.DataChan <- result
		// i++
	}
}

func (ds *DataSourceMySQL) Task() {
	for {
		d := <-ds.DataChan
		sl := storage.Factory("lucene")
		dlen := len(d)
		for i := 0; i < dlen; i++ {
			sl.Add(d[i])
		}

		v := SaveStatus{
			ID:          d[dlen-1]["id"],
			CurrentTime: "1",
		}

		b, err := json.Marshal(v)
		if err != nil {
			fmt.Println("error:", err)
		}

		err = ioutil.WriteFile("conf/test/__tmp.json", b, 0777)
		if err != nil {
			fmt.Println("error:", err)
		}

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
