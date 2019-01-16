package main

import(
	"github.com/newham/hamgo"
	"time"
	"sort"
)

func main(){
	server := hamgo.NewUseConf("app.conf").Server().Static("public")
	server.Get("/", Index)
	server.Get("/info", Info)
	server.Post("/status/connect", Connect)
	server.Post("/status/disconnect", Disconnect)
	server.Get("/status/list", StatusList)
	server.RunAt(hamgo.Conf.String("port"))
}

func Index(ctx hamgo.Context) {
	ctx.HTML("index.html")
}

func Info(ctx hamgo.Context) {
	ctx.HTML("info.html")
}

type Status struct{
	UserAgent string `json:"userAgent"`
	Name string `json:"name"`
	Time int64 `json:"time"`
	IP string `json:"ip"`
	Online bool `json:"online"`
}

var StatusMap = map[string]Status{}

func Connect(ctx hamgo.Context){
	status :=Status{}
	ctx.BindJSON(&status)
	if status.UserAgent==""{
		ctx.WriteString("error")
		ctx.Text(400)
		return
	}
	// cookie := ctx.R().Header["Cookie"]
	// println(cookie[0])
	status.Online = true
	status.Time = time.Now().Unix()
	status.IP=ctx.R().RemoteAddr
	StatusMap[status.UserAgent]=status
	ctx.WriteString("success")
	ctx.Text(200)
}

func Disconnect(ctx hamgo.Context){
	status :=Status{}
	ctx.BindJSON(&status)
	if status.UserAgent!=""{
		delete(StatusMap, status.UserAgent)
	}
	ctx.WriteString("success")
	ctx.Text(200)
}

func StatusList(ctx hamgo.Context){
	statusList :=[]Status{}
	var newMp = make([]string, 0)
   for k, _ := range StatusMap {
      newMp = append(newMp, k)
   }
   sort.Strings(newMp)
   for _, k := range newMp {
	   status:=StatusMap[k]
	now := time.Now().Unix()
	status.Online= (now - status.Time)< hamgo.Conf.DefaultInt64("timeout",10)
	// status.UserAgent = status.UserAgent[0:32]
	statusList =append(statusList,status)
   }
	ctx.JSONFrom(200,&statusList)
}