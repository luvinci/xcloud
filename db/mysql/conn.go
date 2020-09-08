package mysql

import (
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"time"
	cfg "xcloud/config"

	_ "github.com/go-sql-driver/mysql"
)

var (
	user      = cfg.Viper.GetString("mysql.user")
	pwd       = cfg.Viper.GetString("mysql.pwd")
	host      = cfg.Viper.GetString("mysql.host")
	database  = cfg.Viper.GetString("mysql.database")
	charset   = cfg.Viper.GetString("mysql.charset")
	parseTime = cfg.Viper.GetString("mysql.parseTime")
	loc       = cfg.Viper.GetString("mysql.loc")
)

var db *sql.DB

func Conn() *sql.DB {
	return db
}

func init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%s&loc=%s",
		user, pwd, host, database, charset, parseTime, loc)
	db, _ = sql.Open("mysql", dsn)
	// ------ 连接池设置 ------
	// 设置连接池中的最大闲置连接数
	db.SetMaxIdleConns(10)
	// 设置数据库的最大连接数量
	db.SetMaxOpenConns(100)
	// 设置连接的最大可复用时间
	db.SetConnMaxLifetime(time.Hour)
	if err := db.Ping(); err != nil {
		logrus.Errorf("connect to mysql failed, err: %s", err.Error())
		os.Exit(1)
	}
}

// 将查询到的rows转换成[]map
func ParseRows(rows *sql.Rows) []map[string]interface{} {
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}

	record := make(map[string]interface{})
	records := make([]map[string]interface{}, 0)
	for rows.Next() {
		// 将行数据保存到record字典
		err := rows.Scan(scanArgs...)
		if err != nil {
			panic(err)
		}

		for i, col := range values {
			if col != nil {
				record[columns[i]] = col
			}
		}
		records = append(records, record)
	}
	return records
}
