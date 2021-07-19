package hello

import (
	"testing"
)

const TestCount = 3

func TestHello(t *testing.T) {
	want := "Hello, world."
	if got := Hello(); got != want {
		t.Errorf("Hello() = %q, want %q", got, want)
	}
}

func TestChannel(t *testing.T) {
	t.Log("begin")

	// Set up the pipeline.
	c := gen(2, 3)
	out := sq(c)

	// Consume the output.
	//fmt.Println(<-out) // 4
	t.Log(<-out)
	//fmt.Println(<-out) // 9
	t.Log(<-out)

	t.Log("end")
}

func TestSimpleChannel(t *testing.T) {
	out := make(chan int)
	go func() {
		for i := 0; i < TestCount; i++ {
			t.Logf("%d. sending to out...", i)
			out <- i
			t.Logf("%d. sent to out", i)
		}
		t.Log("channel closing")
		close(out)
		t.Log("channel closed")
	}()

	// consumes output
	for i := 0; i < TestCount; i++ {
		t.Logf("%d. receiving from out...", i)
		v, ok := <-out
		t.Logf("%d. receipt from out(v:%v, ok:%v)", i, v, ok)
	}

	// consumes output from closed channel. this code crashes
	//t.Log("receiving from out...")
	//v, ok := <-out
	//t.Logf("receipt from out(v:%v, ok:%v)", v, ok)
}

func TestBufferedChannel(t *testing.T) {
	out := make(chan int, 5)
	go func() {
		for i := 0; i < TestCount; i++ {
			t.Logf("%d. sending to out...", i)
			out <- i
			t.Logf("%d. sent to out", i)
		}
		t.Log("channel closing")
		close(out)
		t.Log("channel closed")
	}()

	// consumes output
	var i int
	for i = 0; i < TestCount; i++ {
		t.Logf("%d. receiving from out...", i)
		v, ok := <-out
		t.Logf("%d. receipt from out(v:%v, ok:%v)", i, v, ok)
	}

	// consumes output from closed channel
	t.Logf("%d. receiving from out...", i)
	v, ok := <-out
	t.Logf("%d. receipt from out(v:%v, ok:%v)", i, v, ok)
}
