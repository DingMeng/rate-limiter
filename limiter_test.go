package rate_limiter

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestNanoNow(t *testing.T) {
	fmt.Println(NanoNow())
}

func TestLimit(t *testing.T){
	rl := New(3,time.Second*5)
	wg := sync.WaitGroup{}
	wg.Add(5)
	for i:=0;i<5;i++ {
		go func(){
			fmt.Println(rl.Limit())
			wg.Done()
		}()
	}
	wg.Wait()
}
