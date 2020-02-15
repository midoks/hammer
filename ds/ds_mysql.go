package ds

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/midoks/hammer/configure"
	"github.com/midoks/hammer/storage"
	"log"
	"strconv"
	"time"
)

type DataSourceMySQL struct {
	Conn     *sql.DB
	DataChan chan map[int]map[string]string
	SS       SaveStatus
	Conf     *configure.Args
}

func (ds *DataSourceMySQL) getPage(p int, s int) (map[int]map[string]string, error) {

	result := make(map[int]map[string]string)
	p = s * p

	err := ds.SS.Read(ds.getTmpFile())
	mSql := fmt.Sprintf("%s where id>%d limit %d offset %d", ds.getQuerySql(), ds.SS.ID, s, p)
	// if err != nil {
	// 	mSql = fmt.Sprintf("%s limit %d offset %d", ds.getQuerySql(), s, p)
	// }

	log.Println(mSql)
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

		result, err := ds.getPage(i, 1000)
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
		sl := storage.OpenStorage(storage.ENGINE_TYPE_LUCENE)
		dlen := len(d)
		for i := 0; i < dlen; i++ {
			sl.Add(d[i])
		}

		updateLastId, err := strconv.ParseInt(d[dlen-1]["id"], 10, 64)
		if err != nil {
			log.Println(err)
		}

		err = ds.SS.Save(updateLastId, ds.getTmpFile())
		if err != nil {
			log.Println(err)
		}
	}
}

func (ds *DataSourceMySQL) Init(conf *configure.Args) {
	ds.DataChan = make(chan map[int]map[string]string, 1000)
	ds.Conf = conf
	ctime := time.Now().Format("2006-01-02 15:04:05")
	ds.SS = SaveStatus{
		ID:          int64(0),
		CurrentTime: ctime,
	}
	ds.InitConn()
}

func (ds *DataSourceMySQL) getDsn() string {

	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", ds.Conf.Conn.User,
		ds.Conf.Conn.Password,
		ds.Conf.Conn.Localhost,
		ds.Conf.Conn.Port,
		ds.Conf.Conn.Db,
		ds.Conf.Conn.Charset)
}

func (ds *DataSourceMySQL) getTmpFile() string {
	return fmt.Sprintf("%s/%s/__tmp.json", ds.Conf.Path, ds.Conf.AppName)
}

func (ds *DataSourceMySQL) getQuerySql() string {
	return fmt.Sprintf("%s", ds.Conf.Sql)
}

func (ds *DataSourceMySQL) InitConn() error {
	dbDSN := ds.getDsn()
	conn, err := sql.Open("mysql", dbDSN)

	if err != nil {
		return err
	}

	ds.Conn = conn
	return nil
}

func (ds *DataSourceMySQL) GetData() (map[int]map[string]string, error) {

	result := make(map[int]map[string]string)

	mSql := fmt.Sprintf("%s limit 1 offset %d", ds.getQuerySql(), 0)
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
