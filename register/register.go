package register

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
	"ztalk_reg/database"
	"ztalk_reg/utils"
)

const key = "006fef6cce9e2900d49f906bef179bf1"

//Reg register
type Reg struct {
	ut *utils.Ut
	db *database.DB
}

type getveryReq struct {
	Cc      string `json:"cc"`
	Phone   string `json:"phone"`
	Scode   string `json:"scode"`
	Lg      string `json:"lg"`
	Imei    string `json:"imei"`
	Imsi    string `json:"imsi"`
	Company string `json:"company"`
	Os      string `json:"os"`
	Model   string `json:"model"`
	Screenh string `json:"screen_h"`
	Screenw string `json:"screen_w"`
	Source  string `json:"source"`
	Sign    string `json:"sign"`
}
type routerReq struct {
	Cc      string `json:"cc"`
	Phone   string `json:"phone"`
	Scode   string `json:"scode"`
	Imei    string `json:"imei"`
	Imsi    string `json:"imsi"`
	Company string `json:"company"`
	Os      string `json:"os"`
	Model   string `json:"model"`
	Screenh string `json:"screen_h"`
	Screenw string `json:"screen_w"`
	Source  string `json:"source"`
	Sign    string `json:"sign"`
}
type registerReq struct {
	Cc         string `json:"cc"`
	Phone      string `json:"phone"`
	Scode      string `json:"scode"`
	Lg         string `json:"lg"`
	VeryCode   string `json:"verycode"`
	Imei       string `json:"imei"`
	Imsi       string `json:"imsi"`
	Company    string `json:"company"`
	Os         string `json:"os"`
	Model      string `json:"model"`
	Screenh    string `json:"screen_h"`
	Screenw    string `json:"screen_w"`
	Source     string `json:"source"`
	Sign       string `json:"sign"`
	Sourceuuid string `json:"sourceuuid"`
}
type getveryRsp struct {
	Code int    `json:"code"`
	Desc string `json:"desc"`
	Wait int    `json:"wait"`
}
type registerRsp struct {
	Code   int    `json:"code"`
	Desc   string `json:"desc"`
	Cc     string `json:"cc"`
	Phone  string `json:"phone"`
	Passwd string `json:"passwd"`
}

type gate struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}
type resources struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}
type routerRsp struct {
	Code    int         `json:"code"`
	Gate    gate        `json:"gate"`
	Reslist []resources `json:"reslist"`
}

//NewRegister new
func NewRegister(db *database.DB, ut *utils.Ut) *Reg {
	return &Reg{
		ut: ut,
		db: db,
	}
}

//Init init
func (r *Reg) Init() {

	http.HandleFunc("/ztalk/getvery", r.getvery)
	http.HandleFunc("/ztalk/register", r.register)
	http.HandleFunc("/ztalk/router", r.router)
	log.Fatal(http.ListenAndServe("0.0.0.0:1101", nil))

}

func (r *Reg) getvery(w http.ResponseWriter, req *http.Request) {

	body, _ := ioutil.ReadAll(req.Body)
	strBody := string(body)
	log.Println(strBody)
	data := getveryReq{}
	ret := getveryRsp{}
	if err := json.Unmarshal(body, &data); err == nil {
		var checksign = data.Phone + data.Imsi + data.Imei + data.Source + key
		if data.Sign != r.ut.Md5(checksign) {
			ret.Code = 0
			ret.Desc = "sign error"
			ret.Wait = 60
		} else {
			ret.Code = 1
			ret.Wait = 60
		}
		rsp, _ := json.Marshal(ret)
		fmt.Fprint(w, string(rsp))
	} else {
		fmt.Println(err)
	}
}

func (r *Reg) register(w http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	strBody := string(body)
	log.Println(strBody)
	data := registerReq{}
	ret := registerRsp{}
	if err := json.Unmarshal(body, &data); err == nil {
		var checksign = data.Phone + data.VeryCode + data.Imsi + data.Imei + data.Source + key
		if data.Sign != r.ut.Md5(checksign) {
			ret.Code = 0
			ret.Desc = "sign error"
		} else {

			var status int
			query := fmt.Sprintf("SELECT fStatus FROM ttestphone WHERE fPhone='%s'", "+8617600113331")
			err := r.db.QueryOne(query).Scan(&status)
			if err == nil {
				//测试号码
				if status == 0 {

					ret.Code = 1
					ret.Cc = data.Cc
					ret.Phone = data.Phone

					password := r.ut.GetPasswd()

					ret.Passwd = r.ut.Base64encode([]byte(password))
					screenH, _ := strconv.Atoi(data.Screenh)
					screenW, _ := strconv.Atoi(data.Screenw)
					source, _ := strconv.Atoi(data.Source)
					insert := fmt.Sprintf("INSERT INTO tuser (fPhone,fCc,fPassword ,fScode,fImei,fImsi,fOs,fCompany ,fModel,fScreenH,fScreenW,fSource,fSourceUuid,fCreateTime,fLastTime) VALUES('%s','%s','%s','%s','%s','%s','%s','%s','%s',%d,%d,%d,'%s',FROM_UNIXTIME(%d),FROM_UNIXTIME(%d))",
						data.Phone, data.Cc, password, data.Scode, data.Imei, data.Imsi, data.Os, data.Company, data.Model, screenH, screenW, source, data.Sourceuuid, time.Now().Unix(), time.Now().Unix())
					//fmt.Println(insert)
					if ok := r.db.UpdateData(insert); ok == false {
						log.Println("already register")
						update := fmt.Sprintf("UPDATE tuser SET fPassword = '%s' ,fLastTime = FROM_UNIXTIME(%d) WHERE fPhone = '%s'", password, time.Now().Unix(), data.Phone)
						ok = r.db.UpdateData(update)
						if ok == false {
							log.Println("register error")
						}
					}
				} else {

				}
			} else {
				//正式号码
			}

		}
		rsp, _ := json.Marshal(ret)
		fmt.Fprint(w, string(rsp))
	} else {
		fmt.Println(err)
	}
}
func (r *Reg) router(w http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	strBody := string(body)
	log.Println(strBody)
	data := routerReq{}
	ret := routerRsp{}
	if err := json.Unmarshal(body, &data); err == nil {
		var checksign = data.Phone + data.Imsi + data.Imei + data.Source + key
		if data.Sign == r.ut.Md5(checksign) {
			ret.Code = 1
			ret.Gate.IP = "192.168.0.98"
			ret.Gate.Port = 8000
			ret.Reslist = []resources{
				resources{IP: "192.168.0.98", Port: 9090},
			}
		} else {
			ret.Code = 0
		}
		rsp, _ := json.Marshal(ret)
		fmt.Fprint(w, string(rsp))
	} else {
		log.Println(err)
	}
}
