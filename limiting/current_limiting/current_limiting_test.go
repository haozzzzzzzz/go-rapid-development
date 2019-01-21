package current_limiting

import (
	"fmt"
	"testing"
	"time"
)

func TestSlice(t *testing.T) {
	a := []int{1, 2, 3, 4}
	fmt.Println(a[:2])
	fmt.Println(a[2:])
}

func TestMemoryStore(t *testing.T) {
	ms := NewMemoryStore()
	items := ms.Get()
	fmt.Println(items)

	ms.Push(1)
	ms.Push(2)
	ms.Push(3)
	ms.Push(4)
	ms.Push(5)
	fmt.Println(ms.Get())

	ms.Pop(1)
	fmt.Println(ms.Get())
	ms.Pop(2)
	fmt.Println(ms.Get())

	ms.Push(6, 7, 8)
	fmt.Println(ms.Get())

	datas := ms.Pop(2)
	fmt.Println(datas)
	fmt.Println(ms.Get())
	datas = append(datas, 9)
	fmt.Println(datas)
	fmt.Println(ms.Get())
}

func TestMinuteFrequencyLimiting_Control(t *testing.T) {
	limiting := &MinuteFrequencyLimiting{
		WorkerNumber: 1,
		Times:        10,
		WaitInterval: time.Second,
		MaxBatchSize: 4,
		Store:        NewMemoryStore(),
		Handler: func(datas ...interface{}) (err error) {
			fmt.Printf("%v %s", datas, time.Now())
			return
		},
	}

	datasC := limiting.Control()
	go func() {
		for datas := range datasC {
			fmt.Printf("%v %s \n", datas, time.Now())
		}
	}()

	for i := 0; i < 100; i++ {
		err := limiting.AcceptData(i)
		if nil != err {
			t.Error(err)
			return
		}

		time.Sleep(600 * time.Millisecond)
	}

	c := make(chan int)
	<-c
}

func TestMinuteFrequencyLimiting_Start(t *testing.T) {
	limiting := &MinuteFrequencyLimiting{
		WorkerNumber: 2,
		Times:        10,
		WaitInterval: time.Second,
		MaxBatchSize: 4,
		Store:        NewMemoryStore(),
		Handler: func(datas ...interface{}) (err error) {
			fmt.Printf("%v %s \n", datas, time.Now())
			return
		},
	}

	limiting.Start()

	for i := 0; i < 100; i++ {
		limiting.AcceptData(i)
		time.Sleep(300 * time.Millisecond)
	}

	c := make(chan int, 0)
	<-c
}
