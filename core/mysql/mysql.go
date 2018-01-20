package mysql

import (
	"GoH/core/log"
	"database/sql"
	"errors"
	"fmt"
	"github.com/dlintw/goconf"
	_ "github.com/go-sql-driver/mysql"
	"regexp"
	"strconv"
	"strings"
	"time"
	"GoH/core/constant"
)

type ORMModel struct {
	tablename string
	params    []string
	column    string
	where     string
	pk        string
	orderby   string
	limit     string
	join      string
}

var (
	DriverName     string
	DataSourceName string
	Orm            *ORMModel
)

func InitConnectionInfo(conf *goconf.ConfigFile) {
	DriverName, _ = conf.GetString("db", "driver_name")
	DataSourceName, _ = conf.GetString("db", "data_source_name")
	if DriverName == "" || DataSourceName == "" {
		fmt.Println("未启用mysql")
		return
	}
}

func OpenMysql() *sql.DB {
	db, err := sql.Open(DriverName, DataSourceName)
	if err != nil {
		panic(err.Error())
	}
	err = db.Ping()
	if err != nil {
		log.Error("Mysql connection failed")
	}
	m := new(ORMModel)
	Orm = m
	return db
}

func (m *ORMModel) FindOne(db *sql.DB) map[int]map[string]interface{} {
	empty := make(map[int]map[string]interface{})
	if db != nil {
		data := m.Limit(1).FindAll(db)
		return data
	}
	log.Error("Mysql connection failed")
	return empty
}

func (m *ORMModel) FindAll(db *sql.DB) map[int]map[string]interface{} {
	result := make(map[int]map[string]interface{})
	if db == nil {
		log.Error("Mysql connection failed")
		return result
	}
	if len(m.params) == 0 {
		m.column = "*"
	} else {
		if len(m.params) == 1 {
			m.column = m.params[0]
		} else {
			m.column = strings.Join(m.params, ",")
		}
	}
	sql := fmt.Sprintf("SELECT %v FROM %v %v %v %v %v", m.column, m.tablename, m.join, m.where, m.orderby, m.limit)
	if constant.Debug {
		log.Debugf("SELECT SQL:%s", sql)
	}
	rows, err := db.Query(sql)
	if err != nil {
		defer func() {
			if err := recover(); err != nil {
				log.Errorf("SQL syntax error, %s:", err)
			}
		}()
		err = errors.New("SQL query error")
	}
	result = QueryResult(rows)
	return result
}

func (m *ORMModel) Insert(db *sql.DB, param map[string]interface{}) (num int, err error) {
	if db == nil {
		return 0, errors.New("Mysql connection failed")
	}
	var keys []string
	var values []string
	if len(m.pk) != 0 {
		delete(param, m.pk)
	}
	for key, value := range param {
		keys = append(keys, key)
		switch value.(type) {
		case int, int32:
			values = append(values, strconv.Itoa(value.(int)))
		case int64:
			values = append(values, strconv.FormatInt(value.(int64), 10))
		case string:
			values = append(values, value.(string))
		case float32, float64:
			values = append(values, strconv.FormatFloat(value.(float64), 'f', -1, 64))
		}
	}
	fileValue := "'" + strings.Join(values, "','") + "'"
	fields := "`" + strings.Join(keys, "`,`") + "`"
	sql := fmt.Sprintf("INSERT INTO %v (%v) VALUES (%v)", m.tablename, fields, fileValue)
	if constant.Debug {
		log.Debugf("INSERT SQL:%s", sql)
	}
	result, err := db.Exec(sql)
	if err != nil {
		defer func() {
			if err := recover(); err != nil {
				log.Errorf("SQL syntax error, %s:", err)
			}
		}()
		err = errors.New("SQL insert error")
		return 0, err
	}
	i, err := result.LastInsertId()
	s, _ := strconv.Atoi(strconv.FormatInt(i, 10))
	if err != nil {
		err = errors.New("insert error")
	}
	return s, err
	
}

func (m *ORMModel) Fileds(param ...string) *ORMModel {
	m.params = param
	return m
}

func (m *ORMModel) Update(db *sql.DB, param map[string]interface{}) (num int, err error) {
	if db == nil {
		return 0, errors.New("Mysql connection failed")
	}
	var setValue []string
	for key, value := range param {
		switch value.(type) {
		case int, int64, int32:
			set := fmt.Sprintf("%v = %v", key, value.(int))
			setValue = append(setValue, set)
		case string:
			set := fmt.Sprintf("%v = '%v'", key, value.(string))
			setValue = append(setValue, set)
		case float32, float64:
			set := fmt.Sprintf("%v = '%v'", key, strconv.FormatFloat(value.(float64), 'f', -1, 64))
			setValue = append(setValue, set)
		}
	}
	setData := strings.Join(setValue, ",")
	sql := fmt.Sprintf("UPDATE %v SET %v %v", m.tablename, setData, m.where)
	if constant.Debug {
		log.Debugf("UPDATE SQL:%s", sql)
	}
	result, err := db.Exec(sql)
	if err != nil {
		defer func() {
			if err := recover(); err != nil {
				log.Errorf("SQL syntax error, %s:", err)
			}
		}()
		err = errors.New("SQL update error")
		return 0, err
	}
	i, err := result.RowsAffected()
	if err != nil {
		err = errors.New("update error")
		return 0, err
	}
	s, _ := strconv.Atoi(strconv.FormatInt(i, 10))
	return s, err
}

func (m *ORMModel) Delete(db *sql.DB, param string) (num int, err error) {
	if db == nil {
		return 0, errors.New("Mysql connection failed")
	}
	h := m.Where(param).FindOne(db)
	if len(h) == 0 {
		return 0, errors.New("The data to be deleted was not found")
	}
	sql := fmt.Sprintf("DELETE FROM %v WHERE %v", m.tablename, param)
	if constant.Debug {
		log.Debugf("DELETE SQL:%s", sql)
	}
	result, err := db.Exec(sql)
	if err != nil {
		defer func() {
			if err := recover(); err != nil {
				log.Errorf("SQL syntax error, %s:", err)
			}
		}()
		err = errors.New("SQL delete error")
		return 0, err
	}
	i, err := result.RowsAffected()
	s, _ := strconv.Atoi(strconv.FormatInt(i, 10))
	if i == 0 {
		err = errors.New("delete error")
	}
	return s, err
}

func (m *ORMModel) Query(db *sql.DB, sql string) interface{} {
	if db == nil {
		return errors.New("Mysql connection failed")
	}
	var query = strings.TrimSpace(sql)
	s, err := regexp.MatchString(`(?i)^select`, query)
	if err == nil && s == true {
		result, _ := db.Query(sql)
		c := QueryResult(result)
		return c
	}
	exec, err := regexp.MatchString(`(?i)^(update|delete)`, query)
	if err == nil && exec == true {
		m_exec, err := db.Exec(query)
		if err != nil {
			return err
		}
		num, _ := m_exec.RowsAffected()
		id := strconv.FormatInt(num, 10)
		return id
	}
	insert, err := regexp.MatchString(`(?i)^insert`, query)
	if err == nil && insert == true {
		m_exec, err := db.Exec(query)
		if err != nil {
			return err
		}
		num, _ := m_exec.LastInsertId()
		id := strconv.FormatInt(num, 10)
		return id
	}
	result, _ := db.Exec(query)
	return result
	
}

func QueryResult(rows *sql.Rows) map[int]map[string]interface{} {
	var result = make(map[int]map[string]interface{})
	columns, _ := rows.Columns()
	columnTypes, _ := rows.ColumnTypes()
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var n = 0
	for rows.Next() {
		result[n] = make(map[string]interface{})
		err := rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
		}
		for i, v := range values {
			result[n][columns[i]] = typeConversion(string(v), columnTypes[i].DatabaseTypeName())
		}
		n++
	}
	return result
}

//类型转换
func typeConversion(value string, ntype string) interface{} {
	var result interface{}
	switch ntype {
	case "INT", "TINYINT":
		result, _ = strconv.ParseInt(value, 10, 32)
	case "BIGINT":
		result, _ = strconv.ParseInt(value, 10, 64)
	case "CHAR", "VARCHAR":
		result = string(value)
	case "DATE", "TIMESTAMP", "DATETIME":
		result, _ = time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
	default:
		result = string(value)
	}
	return result
}

func (m *ORMModel) SetTable(tablename string) *ORMModel {
	m.tablename = tablename
	return m
}

func (m *ORMModel) Where(param string) *ORMModel {
	m.where = fmt.Sprintf("WHERE %v", param)
	return m
}

func (m *ORMModel) SetPk(pk string) *ORMModel {
	m.pk = pk
	return m
}

func (m *ORMModel) OrderBy(param string) *ORMModel {
	m.orderby = fmt.Sprintf("ORDER BY %v", param)
	return m
}

func (m *ORMModel) Limit(size ...int) *ORMModel {
	var end int
	start := size[0]
	if len(size) > 1 {
		end = size[1]
		m.limit = fmt.Sprintf("LIMIT %d,%d", start, end)
		return m
	}
	m.limit = fmt.Sprintf("LIMIT %d", start)
	return m
}

func (m *ORMModel) NoLimit() *ORMModel {
	m.limit = ""
	return m
}

func (m *ORMModel) LeftJoin(table, condition string) *ORMModel {
	m.join = fmt.Sprintf("LEFT JOIN %v ON %v", table, condition)
	return m
}

func (m *ORMModel) RightJoin(table, condition string) *ORMModel {
	m.join = fmt.Sprintf("RIGHT JOIN %v ON %v", table, condition)
	return m
}

func (m *ORMModel) Join(table, condition string) *ORMModel {
	m.join = fmt.Sprintf("INNER JOIN %v ON %v", table, condition)
	return m
}

func (m *ORMModel) FullJoin(table, condition string) *ORMModel {
	m.join = fmt.Sprintf("FULL JOIN %v ON %v", table, condition)
	return m
}
