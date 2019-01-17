package db

import(
	"sync"
	"sort"
	"time"
	"github.com/newham/hamgo"
)

var l sync.Mutex
var StatusMap  map[string]map[string]Status =make(map[string]map[string]Status)

type Status struct{
	UserAgent string `json:"userAgent"`
	Name string `json:"name"`
	Time int64 `json:"time"`
	IP string `json:"ip"`
	Online bool `json:"online"`
	Icon string `json:"icon"`
	OS string `json:"os"`
}

func GetStatus(username ,name string)Status{
	l.Lock()
	defer l.Unlock()
	return StatusMap[username][name]
}

func SetStatus(username,name string,newStatus Status){
	userStatusMap := StatusMap[username]
	if userStatusMap==nil{
		userStatusMap = make(map[string]Status)
	}
	userStatusMap[name]= newStatus
	StatusMap[username] = userStatusMap
}

func DeleteStatus(username,name string){
	
	delete(StatusMap[username],name)
}

func GetStatusList(username string)[]Status{
	statusList :=[]Status{}
	var newMp = make([]string, 0)
	userStatusMap :=StatusMap[username]
	l.Lock()
	defer l.Unlock()
	for k, _ := range userStatusMap {
	   newMp = append(newMp, k)
	}
	
	sort.Strings(newMp)
	
	for _, k := range newMp {
		status:=userStatusMap[k]
	 now := time.Now().Unix()
	 status.Online= (now - status.Time)< hamgo.Conf.DefaultInt64("timeout",10)
	 // status.UserAgent = status.UserAgent[0:32]
	 statusList =append(statusList,status)
	}
	return statusList
}