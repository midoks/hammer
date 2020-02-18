package ds

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/midoks/hammer/configure"
	_ "github.com/midoks/hammer/storage"
	"log"
	"strconv"
	"strings"
	"time"
)

var (
	CONSTANT_PK   = "id"
	CONSTANT_STEP = 1000
)

type DataSourceMySQL struct {
	Conn     *sql.DB
	DataChan chan map[int]map[string]string
	SS       SaveStatus
	Conf     *configure.Args
}

func (ds *DataSourceMySQL) getResult(sql string) (map[int]map[string]string, error) {
	result := make(map[int]map[string]string)

	rows, err := ds.Conn.Query(sql)
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

func (ds *DataSourceMySQL) getPage(p int64, s int) (map[int]map[string]string, error) {

	result := make(map[int]map[string]string)
	// offset := p * int64(s)

	query, err := ds.getQuerySql()
	if err != nil {
		return result, err
	}

	err = ds.SS.Read(ds.getTmpFile())
	mSql := ""

	if err == nil {
		pk := ds.getPk()
		mSql = fmt.Sprintf("%s where %s>%d limit %d", query, pk, ds.SS.PK, s)
	} else {
		mSql = fmt.Sprintf("%s limit %d offset %d", query, s, p)
	}

	return ds.getResult(mSql)
}

// 全量到入
func (ds *DataSourceMySQL) Import() {

	var i int64 = 0
	for {

		result, err := ds.getPage(i, ds.getPageStep())
		if err != nil {
			break
		}

		if len(result) == 0 {
			break
		}

		ds.DataChan <- result
		i++

		//间隔时间
		time.Sleep(1 * time.Second)
	}
}

//增量导入
func (ds *DataSourceMySQL) DeltaData() {

	deltaSql, err := ds.getDeltaQuerySql()

	if err != nil {
		log.Println(deltaSql, err)
		return
	}

	deltaResult, err := ds.getResult(deltaSql)
	if err != nil {
		log.Println(deltaResult, err)
		return
	}

	pk := ds.getPk()
	for i := 0; i < len(deltaResult); i++ {

		pkSql, err := ds.getDeltaImportQuerySql(deltaResult[i][pk])

		if err != nil {
			log.Println(pkSql, err)
			continue
		}

		result, err := ds.getResult(pkSql)
		if err != nil {
			log.Println(result, err)
			continue
		}
		ds.DataChan <- result
	}
}

//删除无效数据
func (ds *DataSourceMySQL) DeleteData() {

}

func (ds *DataSourceMySQL) Task() {
	pk := ds.getPk()
	for {
		d := <-ds.DataChan
		dlen := len(d)
		sl := storage.OpenStorage(storage.ENGINE_TYPE_LUCENE)

		for i := 0; i < dlen; i++ {
			sl.Add(d[i])
		}

		updateLastId, err := strconv.ParseInt(d[dlen-1][pk], 10, 64)
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
		PK:              int64(0),
		LastUpdatedTime: ctime,
	}
	ds.InitConn()
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

func (ds *DataSourceMySQL) getQuerySql() (string, error) {
	if ds.Conf.Query == "" {
		return ds.Conf.Query, errors.New("get query is empty!")
	}
	return ds.Conf.Query, nil
}

func (ds *DataSourceMySQL) getDeltaQuerySql() (string, error) {
	if ds.Conf.DeltaQuery == "" {
		return ds.Conf.DeltaQuery, errors.New("get delta query is empty!")
	}

	err := ds.SS.Read(ds.getTmpFile())
	if err != nil {
		return ds.Conf.DeltaQuery, errors.New("get delta query is fail, program is not running!")
	}

	ds.Conf.DeltaQuery = strings.Replace(ds.Conf.DeltaQuery, "${LAST_UPDATE_TIME}", ds.SS.LastUpdatedTime, -1)
	return ds.Conf.DeltaQuery, nil
}

func (ds *DataSourceMySQL) getDeltaImportQuerySql(pk string) (string, error) {
	if ds.Conf.DeltaImportQuery == "" {
		return ds.Conf.DeltaImportQuery, errors.New("get delta import query is empty!")
	}
	pkSQL := strings.Replace(ds.Conf.DeltaImportQuery, "${PK}", pk, -1)
	return pkSQL, nil
}

func (ds *DataSourceMySQL) getDeletePkQuerySql(pk string) (string, error) {

	if ds.Conf.DeletedPkQuery == "" {
		return ds.Conf.DeletedPkQuery, errors.New("get deleted pk query is empty!")
	}
	pkSQL := strings.Replace(ds.Conf.DeletedPkQuery, "${PK}", pk, -1)
	return pkSQL, nil
}

func (ds *DataSourceMySQL) getDeleteQuerySql() (string, error) {

	if ds.Conf.DeletedQuery == "" {
		return ds.Conf.DeletedQuery, errors.New("get delete query is empty!")
	}
	err := ds.SS.Read(ds.getTmpFile())
	if err != nil {
		return ds.Conf.DeletedQuery, errors.New("get deleted query is fail, program is not running!")
	}

	ds.Conf.DeletedQuery = strings.Replace(ds.Conf.DeletedQuery, "${LAST_UPDATE_TIME}", ds.SS.LastUpdatedTime, -1)
	return ds.Conf.DeletedQuery, nil
}

func (ds *DataSourceMySQL) getPk() string {
	if ds.Conf.Pk == "" {
		return CONSTANT_PK
	}
	return ds.Conf.Pk
}

func (ds *DataSourceMySQL) getPageStep() int {
	if ds.Conf.Step == 0 {
		return CONSTANT_STEP
	}
	return ds.Conf.Step
}
