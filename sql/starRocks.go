// Package sql @author: Violet-Eva @date  : 2025/7/21 @notes :
package sql

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type StarRocks struct {
	*gorm.DB
	sqlDB        *sql.DB
	host         string
	port         int
	feHttpPort   int
	user         string
	password     string
	queryTimeout time.Duration
}

func (s *StarRocks) SetQueryTimeout(queryTimeout int) *StarRocks {
	s.queryTimeout = time.Duration(queryTimeout) * time.Second
	return s
}

func (s *StarRocks) Close() error {
	err := s.sqlDB.Close()
	return err
}

type StarRocksBuilder struct {
	Host          string
	Port          int
	feHttpPort    int
	User          string
	PassWord      string
	dbName        string
	queryTimeout  time.Duration
	retryTimes    int
	retryInterval time.Duration
	maxIdleConn   int
	maxOpenConn   int
	maxLifetime   time.Duration
	opts          []string
	gormConfig    *gorm.Config
}

func NewStarRocksBuilder(host string, port int, user, password string) *StarRocksBuilder {
	return &StarRocksBuilder{
		Host:         host,
		Port:         port,
		User:         user,
		PassWord:     password,
		feHttpPort:   8030,
		queryTimeout: 300 * time.Second,
		maxIdleConn:  2,
		maxOpenConn:  2,
		maxLifetime:  1800 * time.Second,
		opts: []string{
			"timeout=10s",
			"charset=utf8mb4",
		},
	}
}

func (s *StarRocksBuilder) SetFeHttpPort(port int) *StarRocksBuilder {
	s.feHttpPort = port
	return s
}

func (s *StarRocksBuilder) SetDbName(dbName string) *StarRocksBuilder {
	s.dbName = dbName
	return s
}

func (s *StarRocksBuilder) SetQueryTimeout(queryTimeout int) *StarRocksBuilder {
	s.queryTimeout = time.Duration(queryTimeout) * time.Second
	return s
}

func (s *StarRocksBuilder) SetRetry(retryTimes, retryInterval int) *StarRocksBuilder {
	s.retryTimes = retryTimes
	s.retryInterval = time.Duration(retryInterval) * time.Second
	return s
}

func (s *StarRocksBuilder) SetMaxIdleConn(maxIdleConn int) *StarRocksBuilder {
	s.maxIdleConn = maxIdleConn
	return s
}

func (s *StarRocksBuilder) SetMaxOpenConn(maxOpenConn int) *StarRocksBuilder {
	s.maxOpenConn = maxOpenConn
	return s
}

func (s *StarRocksBuilder) SetMaxLifetime(maxLifetime int) *StarRocksBuilder {
	s.maxLifetime = time.Duration(maxLifetime) * time.Second
	return s
}

func (s *StarRocksBuilder) SetOpts(opts []string) *StarRocksBuilder {
	s.opts = opts
	return s
}

func (s *StarRocksBuilder) AddOpt(key, value string) *StarRocksBuilder {
	s.opts = append(s.opts, fmt.Sprintf("%s=%s", key, value))
	return s
}

func (s *StarRocksBuilder) SetGormConfig(config *gorm.Config) *StarRocksBuilder {
	s.gormConfig = config
	return s
}

func (s *StarRocksBuilder) Build() (*StarRocks, error) {
	var (
		opt      string
		mysqlDSN = "%s:%s@tcp(%s:%d)/%s%s"
	)
	if len(s.opts) > 0 {
		opt = "?" + strings.Join(s.opts, "&")
	} else {
		opt = ""
	}
	dsn := fmt.Sprintf(mysqlDSN, s.User, s.PassWord, s.Host, s.Port, s.dbName, opt)
	db, err := gorm.Open(mysql.Open(dsn), s.gormConfig)
	if err != nil {
		return nil, err
	}
	sqlDB, sqlDBErr := db.DB()
	if sqlDBErr != nil {
		return nil, sqlDBErr
	}

	sqlDB.SetMaxIdleConns(s.maxIdleConn)
	sqlDB.SetMaxOpenConns(s.maxOpenConn)
	sqlDB.SetConnMaxLifetime(s.maxLifetime)
	err = sqlDB.Ping()
	if err != nil {
		return nil, err
	}

	return &StarRocks{db, sqlDB, s.Host, s.Port, s.feHttpPort, s.User, s.PassWord, s.queryTimeout}, nil
}

type DescResult struct {
	Field   string `gorm:"column:field"`
	Type    string `gorm:"column:type"`
	Null    string `gorm:"column:null"`
	Key     string `gorm:"column:key"`
	Default string `gorm:"column:default"`
	Extra   string `gorm:"column:extra"`
}
type StreamLoadResult struct {
	TxnId                  int    `json:"TxnId"`
	Label                  string `json:"Label"`
	Status                 string `json:"Status"`
	Message                string `json:"Message"`
	NumberTotalRows        int    `json:"NumberTotalRows"`
	NumberLoadedRows       int    `json:"NumberLoadedRows"`
	NumberFilteredRows     int    `json:"NumberFilteredRows"`
	NumberUnselectedRows   int    `json:"NumberUnselectedRows"`
	LoadBytes              int64  `json:"LoadBytes"`
	LoadTimeMs             int    `json:"LoadTimeMs"`
	BeginTxnTimeMs         int    `json:"BeginTxnTimeMs"`
	StreamLoadPlanTimeMs   int    `json:"StreamLoadPlanTimeMs"`
	ReadDataTimeMs         int    `json:"ReadDataTimeMs"`
	WriteDataTimeMs        int    `json:"WriteDataTimeMs"`
	CommitAndPublishTimeMs int    `json:"CommitAndPublishTimeMs"`
}

func (s *StarRocks) JsonStreamingLoadToStarRocks(dbName, tblName string, body []byte, columns ...string) (slr StreamLoadResult, err error) {
	var (
		dr       []DescResult
		req      *http.Request
		resp     *http.Response
		headers  = make(map[string]string)
		respBody []byte
	)

	if !(len(columns) > 0 && columns[0] != "") {
		descSQL := fmt.Sprintf("desc %s.%s", dbName, tblName)
		err = s.Raw(descSQL).Find(&dr).Error
		if err != nil {
			return
		}

		for _, row := range dr {
			tmpColumn := "\"$." + row.Field + "\""
			columns = append(columns, tmpColumn)
		}
	}

	jsonPaths := "[" + strings.Join(columns, ",") + "]"

	labelName := fmt.Sprintf("%s_%s_%s_%d", dbName, tblName, time.Now().Format("20060102150405"), time.Now().Nanosecond())

	req, err = http.NewRequest("PUT", fmt.Sprintf("http://%s:%d/api/%s/%s/_stream_load", s.host, s.feHttpPort, dbName, tblName), bytes.NewBuffer(body))
	if err != nil {
		return
	}

	headers["label"] = labelName
	headers["Expect"] = "100-continue"
	headers["timezone"] = "Asia/Shanghai"
	headers["format"] = "json"
	headers["max_filter_ratio"] = "0"
	headers["strip_outer_array"] = "true"
	headers["ignore_json_size"] = "true"
	headers["jsonpaths"] = jsonPaths

	for key, val := range headers {
		req.Header.Set(key, val)
	}
	req.SetBasicAuth(s.user, s.password)

	resp, err = (&http.Client{
		Timeout: s.queryTimeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			fmt.Printf("stream load be url is: %s\n", req.URL)
			for key, val := range via[0].Header {
				if key == "Authorization" {
					req.Header[key] = val
				}
			}
			return nil
		},
	}).Do(req)

	if err != nil {
		return
	}

	defer resp.Body.Close()
	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(respBody, &slr)
	if err != nil {
		return
	}
	if strings.ToLower(slr.Status) != "success" {
		err = fmt.Errorf("%s", slr.Message)
	}
	return
}

func (s *StarRocks) ExecQuery(query string, args ...interface{}) (list []map[string]interface{}, err error) {
	defer func() {
		if result := recover(); result != nil {
			err = fmt.Errorf("panic: %v", result)
		}
	}()

	var (
		rows    *sql.Rows
		columns []string
	)

	timeout, cancelFunc := context.WithTimeout(context.Background(), s.queryTimeout)
	defer cancelFunc()
	rows, err = s.sqlDB.QueryContext(timeout, query, args...)
	if err != nil {
		return
	}
	defer rows.Close()
	columns, err = rows.Columns()
	if err != nil {
		return
	}
	rawBytes := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(rawBytes))
	for i := range columns {
		scanArgs[i] = &rawBytes[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return
		}
		record := make(map[string]interface{})
		for index, value := range rawBytes {
			record[columns[index]] = string(value)
		}
		list = append(list, record)
	}
	return list, nil
}

func (s *StarRocks) ExecQueryBatchProcessing(query string, batchSize int, function ...func(input []map[string]interface{}) error) (err error) {
	defer func() {
		if result := recover(); result != nil {
			err = fmt.Errorf("panic: %v", result)
		}
	}()

	var (
		rows    *sql.Rows
		columns []string
	)

	timeout, cancelFunc := context.WithTimeout(context.Background(), s.queryTimeout)
	defer cancelFunc()
	rows, err = s.sqlDB.QueryContext(timeout, query)
	if err != nil {
		return
	}
	defer rows.Close()
	columns, err = rows.Columns()
	if err != nil {
		return
	}
	var list []map[string]interface{}
	rawBytes := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(rawBytes))
	for i := range columns {
		scanArgs[i] = &rawBytes[i]
	}
	rowCount := 0
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return
		}
		record := make(map[string]interface{})
		for index, value := range rawBytes {
			record[columns[index]] = string(value)
		}
		list = append(list, record)
		rowCount++
		if rowCount%batchSize == 0 {
			for _, fun := range function {
				err = fun(list)
				if err != nil {
					return
				}
			}
			list = list[:0]
		}
	}
	if len(list) > 0 {
		for _, fun := range function {
			err = fun(list)
			if err != nil {
				return
			}
		}
	}
	return
}

func (s *StarRocks) ExecQueryToString(query string, args ...interface{}) (list []map[string]string, err error) {
	defer func() {
		if result := recover(); result != nil {
			err = fmt.Errorf("panic: %v", result)
		}
	}()

	var (
		rows    *sql.Rows
		columns []string
	)

	timeout, cancelFunc := context.WithTimeout(context.Background(), s.queryTimeout)
	defer cancelFunc()
	rows, err = s.sqlDB.QueryContext(timeout, query, args...)
	if err != nil {
		return
	}
	defer rows.Close()
	columns, err = rows.Columns()
	if err != nil {
		return
	}
	rawBytes := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(rawBytes))
	for i := range columns {
		scanArgs[i] = &rawBytes[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return
		}
		record := make(map[string]string)
		for index, value := range rawBytes {
			record[columns[index]] = string(value)
		}
		list = append(list, record)
	}
	return
}

func (s *StarRocks) ExecQueryToStruct(query string, data *any, args ...interface{}) (err error) {
	defer func() {
		if result := recover(); result != nil {
			err = fmt.Errorf("panic: %v", result)
		}
	}()

	valueOf := reflect.ValueOf(data)
	if valueOf.Kind() != reflect.Ptr {
		return fmt.Errorf("data is not a pointer")
	}

	var (
		rows    *sql.Rows
		columns []string
		bytes   []byte
	)

	timeout, cancelFunc := context.WithTimeout(context.Background(), s.queryTimeout)
	defer cancelFunc()
	rows, err = s.sqlDB.QueryContext(timeout, query, args...)
	if err != nil {
		return
	}
	defer rows.Close()
	columns, err = rows.Columns()
	if err != nil {
		return
	}
	rawBytes := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(rawBytes))
	for i := range columns {
		scanArgs[i] = &rawBytes[i]
	}

	var list []map[string]string
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return
		}
		record := make(map[string]string)
		for index, value := range rawBytes {
			record[columns[index]] = string(value)
		}
		list = append(list, record)
	}
	bytes, err = json.Marshal(list)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return
	}
	return nil
}

func (s *StarRocks) ExecQueryNoResult(query string, args ...interface{}) (err error) {
	defer func() {
		if result := recover(); result != nil {
			err = fmt.Errorf("panic: %v", result)
		}
	}()

	timeout, cancelFunc := context.WithTimeout(context.Background(), s.queryTimeout)
	defer cancelFunc()
	_, err = s.sqlDB.ExecContext(timeout, query, args...)
	if err != nil {
		return
	}
	return
}
