package timeoutchecker

import (
	"dana-tech.com/pg-device-conn/p2p/libs"
	"dana-tech.com/wbw/logs"
	"dana-tech.com/wbw/util"
	"strings"
	"time"
)

var AccountCache *util.UtilMap

func init() {
	AccountCache = util.NewUtilMap()
}

//
func CheckOffLine() {
	// 1， 从配置文件获得检查周期.
	// 2， 从配置文件获得判断离线的周期.
	// 3， 对于判断为离线的设备修改数据库
	//return
	for {
		var offlinecounts int64
		start := time.Now()
		//
		accounts := AccountCache.Items()
		cur_s := util.GetTimeSec()
		//
		for k, v := range accounts {
			// key , value
			_t := strings.Split(k.(string), "-")
			if len(_t) < 2 {
				continue
			}
			if _t[0] == "U" || v.(*libs.AccountInfo).IsExist == false || v.(*libs.AccountInfo).OffLineFlag == 2 {
				continue
			}
			// Device部分
			if cur_s-v.(*libs.AccountInfo).LastUpdate > int64(48) {
				offlinecounts++
				AccountCache.Delete(k)
				//
				logs.Logger.Warnf("CheckOffLine, device_id: %s, offline.", v.(*libs.AccountInfo).Userid)
			} else {

			}
		}
		logs.Logger.Info("spend time %v  offlinecounts %d", time.Now().Sub(start), offlinecounts)
		logs.Logger.Info("accountcache len %v", AccountCache.Len())
		util.Sleep(10)
		// sleep
		//util.Sleep(58 - 10)

	}
}
