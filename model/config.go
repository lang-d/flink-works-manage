package model

/*
*
数据库配置
*/
type DataBaseConfig struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	Addr         string `json:"addr"`
	Port         int    `json:"port"`
	DbName       string `json:"dbName"`
	MaxOpenConns int    `json:"maxOpenConns,omitempty"`
	MaxIdeConns  int    `json:"maxIdeConns,omitempty"`
}
