package core

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/land-d/flink-works-amage/model"
)

const DRIVER_NAME = "mysql"

type Results []map[string]interface{}

type Result map[string]interface{}

type DataBase struct {
	pool         *sql.DB
	maxOpenConns int
	maxIdeConns  int
}

type Mysql struct {
	db map[string]*DataBase
}

func NewMysql() *Mysql {
	return &Mysql{db: map[string]*DataBase{}}
}

/*
*
注册一个数据库连接
*/
func (this *Mysql) regiestDb(config model.DataBaseConfig, name string) {
	db := NewDb(config.Username, config.Password, config.Addr, config.Port, config.DbName, config.MaxOpenConns, config.MaxIdeConns)
	this.db[name] = db
}

func NewDb(user string, passwd string, addr string, port int, dbName string, maxOpenConns int, maxIdeConns int) *DataBase {
	db := &DataBase{maxIdeConns: maxIdeConns, maxOpenConns: maxOpenConns}
	return db.init(user, passwd, addr, port, dbName)
}

func (this *DataBase) init(user string, passwd string, addr string, port int, dbName string) *DataBase {
	// use dsn like this username:password@protocol(address)/dbname?param=value
	pool, err := sql.Open(DRIVER_NAME, fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4", user, passwd, addr, port, dbName))
	if err != nil {
		panic(err)
	}
	this.pool = pool
	this.pool.SetMaxOpenConns(this.maxOpenConns)
	this.pool.SetMaxIdleConns(this.maxIdeConns)

	err = this.pool.Ping()
	if err != nil {
		panic(err)
	}
	return this

}

// query one row
func (this *DataBase) QueryOne(sql string, params ...interface{}) Result {
	rows, err := this.pool.Query(sql, params...)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	cloums, err := rows.Columns()
	if err != nil {
		panic(err)
	}
	fields := this.createFeilds(cloums)
	for rows.Next() {
		err = rows.Scan(fields...)
		if err != nil {
			panic(err)
		}
	}
	return this.createResult(cloums, fields)
}

// query all
func (this *DataBase) Query(sql string, params ...interface{}) Results {
	rows, err := this.pool.Query(sql, params...)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	cloums, err := rows.Columns()
	if err != nil {
		panic(err)
	}
	results := make(Results, 0)
	for rows.Next() {
		fields := this.createFeilds(cloums)
		err = rows.Scan(fields...)
		if err != nil {
			panic(err)
		}
		results = append(results, this.createResult(cloums, fields))
	}
	return results
}

/*
*
执行query
*/
func (this *DataBase) Exec(query string, params ...any) (sql.Result, error) {
	return this.pool.Exec(query, params...)
}

/*
*
执行query,先预编译
*/
func (this *DataBase) ExecWithPrepare(query string, params ...any) (sql.Result, error) {
	stmt, err := this.pool.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(params)
}

/*
*
执行事务
*/
func (this *DataBase) ExecWithTx(sqlList []string, params [][]interface{}) (error, error) {

	ctx := context.Background()
	tx, err := this.pool.BeginTx(ctx, &sql.TxOptions{}) //开启事务
	if err != nil {
		return err, nil
	}

	for i := 0; i < len(sqlList); i++ {
		if params[i] == nil || len(params[i]) == 0 {
			_, err = tx.Exec(sqlList[i])
		} else {
			_, err = tx.Exec(sqlList[i], params[i]...)
		}
		if err != nil {
			errRollback := tx.Rollback()
			return err, errRollback
		}

	}
	errCommit := tx.Commit() // 事务提交
	if errCommit != nil {
		errRollback := tx.Rollback()
		return errCommit, errRollback
	}

	return nil, nil

}

func (this *DataBase) createFeilds(cloums []string) []interface{} {
	slice := make([]interface{}, len(cloums))
	for i, _ := range cloums {
		slice[i] = new(interface{})
	}
	return slice
}

func (this *DataBase) createResult(cloums []string, fields []interface{}) Result {
	result := make(Result)
	for i, k := range cloums {
		result[k] = *(fields[i].(*interface{}))
	}
	return result
}
