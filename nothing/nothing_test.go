package nothing

import (
	"testing"
)

func TestPop(t *testing.T) {

	var manage TimerManager
	var node HeapNode
	node.Time = 1
	var node1 HeapNode
	node1.Time = 2
	manage.Push(&node1)
	manage.Push(&node)

	min := manage.Pop()
	if min == nil {
		t.Failed()
	}
	n := min.(*HeapNode)
	t.Log(n.Time)

}
