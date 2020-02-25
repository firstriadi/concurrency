package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

type tambahOp struct {
	a        int
	b        int
	hasil    int
	workerID int
}

func main() {
	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, os.Interrupt)

	exitCuy := make(chan struct{})
	go func() {
		<-interruptCh
		close(exitCuy)
	}()

	chWorker := make(chan *tambahOp)
	chSink := make(chan *tambahOp)

	for i := 0; i < 10; i++ {
		go worker(exitCuy, i, chWorker, chSink)
	}

	go func() {
		// producer
		for {
			tops := tambahOp{a: getRandomInt(), b: getRandomInt()}
			select {
			case chWorker <- &tops:
			case <-exitCuy:
				return
			}
		}
	}()

	// printer
	for {
		select {
		case a := <-chSink:
			fmt.Printf("%d tambah %d adalah %d dikerjain oleh %d\n", a.a, a.b, a.hasil, a.workerID)
		case <-exitCuy:
			return
		}
	}
}

func getRandomInt() int {
	return rand.Int() % 1000
}

func worker(exitCh <-chan struct{}, id int, work <-chan *tambahOp, output chan<- *tambahOp) {
	for w := range work {
		w.hasil = hardTambah(w.a, w.b)
		w.workerID = id
		select {
		case output <- w:
		case <-exitCh:
			return
		}

	}
}

func hardTambah(a, b int) int {
	time.Sleep(0 +
		(500 * time.Millisecond) +
		(time.Duration(rand.Int()%500) * time.Millisecond),
	)
	return a + b
}
