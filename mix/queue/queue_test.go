package queue

import (
	"sort"
	"testing"
)

func TestRunRandomQueue(t *testing.T) {
	window := uint64(3600)    // 1h
	granularity := uint64(10) // 10 seconds. 360 queues
	messagesPerWindow := uint64(3600)
	// -------------
	queue, err := New(window, granularity, messagesPerWindow)
	if err != nil {
		t.Fatalf("New RandomQueue: %s", err)
	}
	if queue.rounds != 360 {
		t.Error("Wrong rounds")
	}
	if len(queue.slices) != 360 {
		t.Error("Wrong number of timeslices")
	}
	if uint64(cap(queue.slices[0].messages)) != 10 {
		t.Errorf("Wrong number of message elements: %d", cap(queue.slices[0].messages))
	}
	for i := uint64(0); i < 7200; i++ {
		queue.Add(3600+i, i, 3600+i)
	}
	res := make([]interface{}, 0, 7203)
	for i := uint64(3600); i < 10000; i += 10 {
		res = append(res, queue.GetSendQueue(i)...)
	}
	if len(res) != 7200 {
		t.Error("Messages lost or duplicate")
	}

	for i := uint64(10000); i < 100000; i += 10 {
		if len(queue.GetSendQueue(i)) > 0 {
			t.Error("GetSendQueue returned elements from empty queue")
		}
	}

	resI := make([]int, len(res))
	for i, e := range res {
		resI[i] = int(e.(uint64))
	}
	sort.Sort(sort.IntSlice(resI))
	last := -1
	for i, j := range resI {
		if j != last+1 {
			t.Fatalf("Missed or extra element: %d, %d", i, last)
		}
		last = j
	}
}
