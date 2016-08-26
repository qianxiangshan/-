package timeoutchecker

import (
	//"math/rand"
	"testing"
	"time"
)

func tt(tmp interface{}) {

}

type Test struct {
}

func BenchmarkUpdatetime(b *testing.B) {

	timeout := OfflineMangeInit(24, tt)

	for i := 0; i < b.N; i++ {
		value := new(Test)
		//timeout.UpdateTimer(0, time.Now().Unix()+rand.Int63()%1, value)
		timeout.UpdateTimer(0, time.Now().Unix(), value)
	}

}

func TestUpdatetime(t *testing.T) {

	timeout := OfflineMangeInit(24, tt)
	value := new(Test)
	lasttime := time.Now().Unix() + 4
	tmp, err := timeout.UpdateTimer(0, lasttime, value)

	t.Log("tmp returrn ", tmp, err)
	timeout.ReportStatus()

	time.Sleep(time.Second)
	nowtime := time.Now().Unix() + 4
	tmp, err = timeout.UpdateTimer(lasttime, nowtime, tmp)
	timeout.ReportStatus()
	t.Log("sec ", tmp, err)

}
