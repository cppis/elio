package hello

func Hello() string {
	return "Hello, world."
}

// Go Concurrency Patterns: Pipelines and cancellation(https://blog.golang.org/pipelines)
// gen pipeline producer
func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

// gen pipeline consumer and producer
func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}
