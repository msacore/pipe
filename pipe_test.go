package pipe

import (
	"fmt"
	"testing"
)

func TestPipe(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		// p := New[string]()
		// pw := p.Writer()
		// pr := p.Reader()
		// go func() {
		// 	pw.Write("YAY!")
		// }()
		// t.Error(pr.Read())

		// ch1 := make(chan int)
		// go func() {
		// 	for i := 0; i < 10; i++ {
		// 		ch1 <- i
		// 	}
		// 	close(ch1)
		// }()
		// ch2 := make(chan int)
		// Filter(ch1, ch2, func(in int) bool { return in%2 == 0 })
		// ch3 := Map(ch2, func(in int) string { return fmt.Sprintf("%d IS OK!", in) })
		// t.Error(<-ch3)
		// t.Error(<-ch3)

		input := make(chan int, 5)
		go func() {
			for i := 0; i < 10; i++ {
				input <- i
			}
			close(input)
		}()
		// filtered := FilterSequential(input, func(value int) bool { return value%2 == 0 })
		filtered := Filter(input, func(value int) bool { return value%2 == 0 })

		// mapped := MapSequential(filtered, func(value int) string { return fmt.Sprintf("%d IS OK!", value) })
		mapped := Map(filtered, func(value int) string { return fmt.Sprintf("%d IS OK!", value) })

		mapped = Watch(mapped, func(value string) { t.Error(value) })
		Wait(mapped)
		// time.Sleep(time.Second * 1)
	})
}
