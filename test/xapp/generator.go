package main

import (
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"
	"sync"
	"time"
)

var (
	wg     sync.WaitGroup
	mux    sync.Mutex
	rx     int
	tx     int
	failed int
)

type Generator struct {
}

func (m Generator) Consume(mtype, subId, len int, payload []byte) (err error) {
	xapp.Logger.Debug("message received - type=%d subId=%d len=%d", mtype, subId, len)

	mux.Lock()
	rx++
	mux.Unlock()

	ack := xapp.Config.GetInt("test.waitForAck")
	if ack != 0 {
		wg.Done()
	}

	return nil
}

func waitForMessages() {
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	// All done!
	case <-time.After(5000 * time.Millisecond):
		xapp.Logger.Warn("Message waiting timed out!")
	}
}

func runTests(mtype, subId, amount, msize, ack int) {
	tx = 0
	rx = 0
	s := make([]byte, msize, msize)

	start := time.Now()
	for i := 0; i < amount; i++ {
		if ok := xapp.Rmr.Send(mtype, subId, msize, s); ok {
			tx++
			if ack != 0 {
				wg.Add(1)
			}
		} else {
			failed++
		}
	}

	// Wait until all replies are received, or timeout occurs
	waitForMessages()

	elapsed := time.Since(start)
	xapp.Logger.Info("amount=%d|tx=%d|rx=%d|failed=%d|time=%v\n", amount, tx, rx, failed, elapsed)
}

func generator() {
	// Start RMR and wait until engine is ready
	go xapp.Rmr.Start(Generator{})
	for xapp.Rmr.IsReady() == false {
		time.Sleep(time.Duration(2) * time.Second)
	}

	// Read parameters
	interval := 1000000 * 1.0 / xapp.Config.GetInt("test.rate")
	mtype := xapp.Config.GetInt("test.mtype")
	subId := xapp.Config.GetInt("test.subId")
	amount := xapp.Config.GetInt("test.amount")
	size := xapp.Config.GetInt("test.size")
	ack := xapp.Config.GetInt("test.waitForAck")
	rounds := xapp.Config.GetInt("test.rounds")

	// Now generate message load as per request
	for i := 0; i < rounds; i++ {
		runTests(mtype, subId, amount, size, ack)
		if interval != 0 {
			time.Sleep(time.Duration(interval) * time.Microsecond)
		}
	}

	return
}
