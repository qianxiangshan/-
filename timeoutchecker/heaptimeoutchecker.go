package timeoutchecker

import (
	"container/heap"
	"container/list"
	"errors"
	"log"
	"sync"
	"time"
	"tool/timermanage/nothing"
)

type OfflineManager struct {
	//heap的实现
	Offlinetimer nothing.TimerManager
	//并发保护heap的锁
	Locker sync.RWMutex
	//超时处理函数
	TimeoutHandler nothing.TimeoutHandler
	//超时间隔
	HbInterval int64
	//时间到链表的映射
	Time2list map[int64]*nothing.HeapNode
}

//init before using
func OfflineMangeInit(hbinterval int64, handler nothing.TimeoutHandler) *OfflineManager {
	OfflineMange := new(OfflineManager)
	//超时处理函数
	OfflineMange.TimeoutHandler = handler
	//超时间隔
	OfflineMange.HbInterval = hbinterval
	//初始化堆
	heap.Init(&OfflineMange.Offlinetimer)
	//;初始化map
	OfflineMange.Time2list = make(map[int64]*nothing.HeapNode)

	return OfflineMange
}

func (OfflineManage *OfflineManager) HeapTimeout() error {

	var currenttime int64
	var lastnode *nothing.HeapNode
	var lasttimeout int64
	var lastnodelistlen int

	if OfflineManage == nil {
		log.Fatal("not init offline manager")
		return errors.New("not init")
	}
	//当前时间距离最近的
	lasttimeout = OfflineManage.HbInterval

	for {

		log.Println("start timeout sleep time ", lasttimeout)
		//第一个心跳周期不用处理
		time.Sleep(time.Second * time.Duration(lasttimeout))
		start := time.Now()
		//最小堆中无数据不用处理
		if OfflineManage.Offlinetimer.Len() == 0 {
			lasttimeout = 5
			continue
		}
		//取得最近的超时时间的数据,即第0个数据.
		OfflineManage.Locker.RLock()
		lastnode = OfflineManage.Offlinetimer[0]
		OfflineManage.Locker.RUnlock()

		//取得当前时间
		currenttime = time.Now().Unix()
		log.Printf("lastnode time %d current %d\n", lastnode.Time, currenttime)
		if lastnode.Time >= currenttime {
			//还没有超时的数据
			lasttimeout = lastnode.Time - currenttime + 1
			continue
		} else {
			//加锁对链表进行删除操作
			lastnode.Locker.Lock()
			lastnodelistlen = lastnode.ListHeader.Len()
			//处理当前超时的数据.
			data := lastnode.ListHeader.Front()

			for data != nil {
				OfflineManage.TimeoutHandler(data.Value)
				tmp := data.Next()
				lastnode.ListHeader.Remove(data)
				data = tmp
			}
			lastnode.Locker.Unlock()
			//删除map中的数据,由于是写操作,所以要加锁防止崩溃
			OfflineManage.Locker.Lock()
			delete(OfflineManage.Time2list, lastnode.Time)
			heap.Pop(&OfflineManage.Offlinetimer)
			OfflineManage.Locker.Unlock()
		}
		//处理完成后需要进行1s的休眠,然后进入下一轮的处理.
		lasttimeout = 0
		log.Printf("handle  %d  spend time %v\n", lastnodelistlen, time.Now().Sub(start))
	}
}

//updatetime and lasttime must be the same with last operation,lasttime=0 means first time to be managed
//该函数不允许在超时处理函数中使用,element 第一次update时使用需要记录的指针,第二次后都是用该函数返回值.
//updatetime 和lasttime都是unix时间戳
//返回值interface,是用户更新来的凭证,因此用户需要保存
func (OfflineManage *OfflineManager) UpdateTimer(lasttime, updatetime int64, element interface{}) (interface{}, error) {
	if element == nil {
		return nil, errors.New("nil pointer")
	}

	if updatetime < time.Now().Unix() {
		return nil, errors.New("args invalid")
	}

	var elementtype bool
	var listelement *list.Element

	OfflineManage.Locker.RLock()
	timernode, has := OfflineManage.Time2list[lasttime]
	OfflineManage.Locker.RUnlock()

	if !has && lasttime != 0 {
		//do nothing
		log.Println("lasttime not int heap")
	} else if lasttime != 0 {
		timernode.Locker.Lock()
		//log.Println("remove old data")
		timernode.ListHeader.Remove(element.(*list.Element))
		timernode.Locker.Unlock()
	} else if lasttime == 0 {
		//新插入数据,不需要移除旧的
		elementtype = true
	}

	//可能插入新的序列,先加读锁控制,防止全部的数据都需要加写锁
	OfflineManage.Locker.RLock()
	timernode, has = OfflineManage.Time2list[updatetime]
	OfflineManage.Locker.RUnlock()

	if !has {
		//log.Printf("timestamp %d not exit\n", updatetime)
		//加写锁判断后在加入防止重复添加
		OfflineManage.Locker.Lock()
		timernode, has = OfflineManage.Time2list[updatetime]
		if !has {
			//没有这个时间戳记录,则增加,加写锁啊
			heapnode := new(nothing.HeapNode)
			heapnode.Time = updatetime
			OfflineManage.Time2list[updatetime] = heapnode
			timernode = heapnode
			heap.Push(&OfflineManage.Offlinetimer, heapnode)
		}
		OfflineManage.Locker.Unlock()
		//log.Printf("timestamp %d insert\n", updatetime)

	}
	if elementtype {
		timernode.Locker.Lock()
		listelement = timernode.ListHeader.PushBack(element)
		timernode.Locker.Unlock()
	} else {
		timernode.Locker.Lock()
		listelement = timernode.ListHeader.PushBack(element.(*list.Element).Value)
		timernode.Locker.Unlock()
	}
	return listelement, nil
}

func (OfflineManage *OfflineManager) ReportStatus() {
	OfflineManage.Locker.RLock()
	l := OfflineManage.Offlinetimer.Len()
	log.Println("report heap size ", l)
	for i := 0; i < l; i++ {
		log.Printf("report timestamp %d, list len %d\n", OfflineManage.Offlinetimer[i].Time, OfflineManage.Offlinetimer[i].ListHeader.Len())
	}
	OfflineManage.Locker.RUnlock()
}
