package completer

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
	"unsafe"
)

func TestCorrectness1(t *testing.T) {
	t.Parallel()
	number := time.Now().UnixNano()
	wg := sync.WaitGroup{}
	wg.Add(3)
	t.Log("number:", number)
	wait, complete := NewCompleter[int64]()
	for range make([]struct{}, 3) {
		go func(t *testing.T) {
			defer wg.Done()
			if n := wait(); n != number {
				panic(fmt.Errorf("expected %d got %d", number, n))
			}
		}(t)
	}
	runtime.Gosched()
	complete(number)
	wg.Wait()
}

func TestCorrectness2(t *testing.T) {
	t.Parallel()
	number := time.Now().UnixNano()
	wg := sync.WaitGroup{}
	wg.Add(3)
	t.Log("number:", number)
	wait, complete := NewCompleter[int64]()
	complete(number)
	runtime.Gosched()
	for range make([]struct{}, 3) {
		go func(t *testing.T) {
			defer wg.Done()
			if n := wait(); n != number {
				panic(fmt.Errorf("expected %d got %d", number, n))
			}
		}(t)
	}
	wg.Wait()
}

func TestCorrectness3(t *testing.T) {
	number := time.Now().UnixNano()
	wait, complete := NewCompleter[int64]()
	complete(number)
	if n := wait(); n != number {
		t.Fatalf("expected %d got %d", number, n)
	}
}
func TestCorrectness4(t *testing.T) {
	wait, complete := NewCompleter[string]()
	func(c1 CompleteFunc[string]) {
		defer func(c2 CompleteFunc[string]) {
			c2("foobar")
		}(c1)
	}(complete)
	if s := wait(); s != "foobar" {
		t.Fatalf("expected 'foobar' got %s", s)
	}
}

func TestCorrectness5(t *testing.T) {
	t.Parallel()
	wait, _ := NewCompleter[string]()

	timer := time.NewTimer(2 * time.Second)
	select {
	case <-func() chan string {
		c := make(chan string, 1)
		go func() {
			c <- wait()
		}()
		return c
	}():
		t.FailNow()
	case <-timer.C:
	}
}
func TestCorrectness6(t *testing.T) {
	t.Parallel()
	wait, complete := NewCompleter[string]()
	timer := time.NewTimer(2 * time.Second)
	go func() {
		time.Sleep(time.Second)
		complete("foobar")
	}()
	select {
	case s := <-func() chan string {
		c := make(chan string, 1)
		go func() {
			c <- wait()
		}()
		return c
	}():
		if s != "foobar" {
			t.FailNow()
		}
	case <-timer.C:
		t.FailNow()
	}
}

func TestCorrectness7(t *testing.T) {
	wait, complete := NewCompleter[string]()
	func(c1 unsafe.Pointer) {
		defer func(c2 uintptr) {
			(*(*CompleteFunc[string])(unsafe.Pointer(c2)))("foobar")
		}(uintptr(c1))
	}(unsafe.Pointer(&complete))
	if s := wait(); s != "foobar" {
		t.Fatalf("expected 'foobar' got %s", s)
	}
}
