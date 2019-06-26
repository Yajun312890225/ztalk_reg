package main

import (
	"ztalk_reg/database"
	"ztalk_reg/register"
	"ztalk_reg/utils"

	_ "github.com/go-sql-driver/mysql"
)

var db = database.NewDB()
var ut = utils.NewUtils()
var redisConn = database.NewRedis()
var reg = register.NewRegister(db, ut, redisConn)

func main() {

	reg.Init()
}
