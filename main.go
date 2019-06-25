package main

import (
	"ztalk_reg/database"
	"ztalk_reg/register"
	"ztalk_reg/utils"

	_ "github.com/go-sql-driver/mysql"
)

var db = database.NewDB()
var ut = utils.NewUtils()
var reg = register.NewRegister(db, ut)

func main() {

	reg.Init()
}
