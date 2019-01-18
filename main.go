package main

import(
	"github.com/newham/hamgo"
	"time"
	"strings"
	"online-checker/db"
)

func main(){
	server := hamgo.NewUseConf("app.conf").UseSession(3600).Server().Static("public")
	server.Filter(SessionFilter).AddAnnoURL("/signin").AddAnnoURL("/favicon.ico").AddAnnoURL("/public")
	server.Get("/favicon.ico",Favicon)
	server.Get("/", Index)
	server.Get("/help", Help)
	server.Handler("/signin", Signin,"POST,GET")
	server.Post("/signout", Signout)
	server.Get("/info", Info)
	server.Post("/status/connect", Connect)
	server.Post("/status/disconnect", Disconnect)
	server.Get("/status/list", StatusList)
	server.RunAt(hamgo.Conf.String("port"))
}

func Favicon(ctx hamgo.Context){
	ctx.File("public/favicon.ico")
}

func SessionFilter(ctx hamgo.Context)bool{
	session :=ctx.GetSession().Get(SESSION_NAME)
	if session!=nil{
		return true
	}
	
	url :=ctx.R().URL.Path
	println("unsignined",url)
	if url == "/"{
		ctx.Redirect("/signin")
	}else{
		ctx.JSONFrom(403,map[string]string{"error":"unsignined"})
	}
	return false
}

func getUsername(ctx hamgo.Context)string{
	signinForm :=ctx.GetSession().Get(SESSION_NAME).(SigninForm)
	return signinForm.Username
}

func Index(ctx hamgo.Context) {
	println("index")
	ctx.PutData("username",getUsername(ctx))
	ctx.HTML("pages/index.html")
}

func Help(ctx hamgo.Context) {
	println("help")
	ctx.HTML("pages/help.html")
}

type SigninForm struct{
	Username string `form:"username"`
	Password string `form:"password"`
}

const SESSION_NAME = "go_session"

func Signin(ctx hamgo.Context) {
	if strings.ToLower(ctx.Method()) == "get"{
		ctx.HTML("pages/signin.html")
	}else {
		signinForm :=SigninForm{}
		ctx.BindForm(&signinForm)
		if (signinForm.Username=="liuhan" && signinForm.Password=="123") || (signinForm.Username=="admin" && signinForm.Password=="123"){
			ctx.GetSession().Set(SESSION_NAME,signinForm)
			ctx.Redirect("/")
			return
		}
		ctx.HTML("pages/signin.html")
	}
	
}

func Signout(ctx hamgo.Context){
	println("signout")
	ctx.GetSession().Delete(SESSION_NAME)
	ctx.Redirect("/signin")
}

func Info(ctx hamgo.Context) {
	ctx.HTML("pages/info.html")
}

func Connect(ctx hamgo.Context){
	status :=db.Status{}
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
	ip:=ctx.R().RemoteAddr
	status.IP =ip[:strings.LastIndex(ip, ":")]
	if(status.IP == "[::1]"){
		status.IP = "127.0.0.1"
	}
	if ICON_MAP[status.OS]!=""{
		status.Icon = ICON_MAP[status.OS]
	}else{
		status.Icon = ICON_MAP["others"]
	}
	
	db.SetStatus(getUsername(ctx),status.UserAgent,status)
	ctx.WriteString("success")
	ctx.Text(200)
}

func Disconnect(ctx hamgo.Context){
	status :=db.Status{}
	ctx.BindJSON(&status)
	if status.UserAgent!=""{
		db.DeleteStatus(getUsername(ctx), status.UserAgent)
	}
	ctx.WriteString("success")
	ctx.Text(200)
}

const ICON_FOLDER = "/public/img/"
var ICON_MAP =map[string]string {
	"windows 10":ICON_FOLDER+"win10.png",
	"windows 7":ICON_FOLDER+"win7.png",
	"mac os x":ICON_FOLDER+"imac-new.png",
	"android":ICON_FOLDER+"android.png",
	"iphone":ICON_FOLDER+"iphone.png",
	"ipad":ICON_FOLDER+"ipad.png",
	"others":ICON_FOLDER+"other.png",
}

func StatusList(ctx hamgo.Context){
	statusList :=db.GetStatusList(getUsername(ctx))
	ctx.JSONFrom(200,&statusList)
}