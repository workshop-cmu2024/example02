package main

import (
	"encoding/binary"
	"flag"
	"log"
	"math/rand"
	"net/netip"
	"runtime"
	dbg "runtime/debug"
	"time"

	"github.com/talostrading/sonic"
	"github.com/talostrading/sonic/multicast"
	"github.com/talostrading/sonic/sonicerrors"
)

var (
	addr  = flag.String("addr", "224.0.0.224:8080", "multicast group address")
	iter  = flag.Int("iter", 10, "how many iterations, if 0, infinite")
	debug = flag.Bool("debug", false, "if true you can see what is sent")
	rate  = flag.Duration("rate", 50*time.Microsecond, "sending rate")

	letters = []byte("abcdefghijklmnopqrstuvwxyz")
)

func main() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	dbg.SetGCPercent(-1) // turn GC off

	flag.Parse()

	maddr, err := netip.ParseAddrPort(*addr)
	if err != nil {
		panic(err)
	}

	ioc := sonic.MustIO()
	defer ioc.Close()

	p, err := multicast.NewUDPPeer(ioc, "udp", "")
	if err != nil {
		panic(err)
	}

	b := make([]byte, 256)
	var (
		current       uint32 = 1
		last          uint32 = 1
		backfill             = false
		backfillUntil uint32 = 1

		// 1 then increment random, backfill is true
		// if backfill is true, we need to backfill until we increment and then
		// is false
	)
	prepare := func() {
		binary.BigEndian.PutUint32(b, current) // sequence number

		// payload size
		n := rand.Intn(len(b) - 8 /* 4 for seq, 4 for payload size */)
		binary.BigEndian.PutUint32(b[4:], uint32(n))

		// payload
		for i := 0; i < n; i++ {
			b[8+i] = letters[i%len(letters)]
		}

		if *debug {
			log.Printf(
				"sending seq=%d n=%d payload=%s",
				current,
				n,
				string(b[8:8+n]),
			)
		}

		if backfill {
			last++
			if last == backfillUntil {
				backfill = false
				current = last + 1
			} else {
				current = last
			}
		} else {
			last = current
			increment := uint32(rand.Int31n(10) + 1)
			current += increment
			backfill = true
			backfillUntil = current
		}

	}

	if *rate > 0 {
		log.Printf("sending at a rate of %s", *rate)

		t, err := sonic.NewTimer(ioc)
		if err != nil {
			panic(err)
		}
		err = t.ScheduleRepeating(*rate, func() {
			prepare()
			// Send the same packet 10 times to increase the likelihood of
			// arrival.
			for i := 0; i < 10; i++ {
				_, err := p.Write(b, maddr)
				for err == sonicerrors.ErrWouldBlock {
					_, err = p.Write(b, maddr)
				}
				if err != nil {
					panic(err)
				}
			}
		})
		if err != nil {
			panic(err)
		}
	} else {
		log.Print("sending as fast as possible")

		var onWrite func(error, int)
		onWrite = func(err error, _ int) {
			if err == nil {
				prepare()
				// Send the same packet 10 times to increase the likelihood of
				// arrival.
				for i := 0; i < 10; i++ {
					p.AsyncWrite(b, maddr, onWrite)
				}
			} else {
				panic(err)
			}
		}

		prepare()
		// Send the same packet 10 times to increase the likelihood of arrival.
		for i := 0; i < 10; i++ {
			p.AsyncWrite(b, maddr, onWrite)
		}
	}

	log.Print("starting...")
	if *iter == 0 {
		for {
			_, _ = ioc.PollOne()
		}
	} else {
		for i := 0; i < *iter; i++ {
			_, _ = ioc.PollOne()
		}
	}
}