package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

var localZone *time.Location

type config struct {
	startAt    int64
	finishAt   int64
	msisdnFile string
	rpd        int
}

func main() {
	localZone, _ = time.LoadLocation("Europe/Prague")
	cfg := mustCreateConfig()
	msisdns := mustLoadMSISDNs(cfg.msisdnFile, cfg.startAt)
	rand.Seed(time.Now().UnixNano())

	yw := newYateWriter(".", "voice-artificial-yate-tsv-")

	var err error

GEN:
	for d := cfg.startAt; d < cfg.finishAt; d = d + 86400 {
		num := cfg.rpd
		tw := time.Unix(d, 0).Weekday()
		if tw == time.Sunday || tw == time.Saturday {
			num /= 10
		}
		recs := day(d, num, msisdns)

		for _, hc := range recs {
			for _, cs := range hc {
				if err := yw.Write(cs); err != nil {
					break GEN
				}
			}
		}

	}
	yw.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
	}

}

func mustLoadMSISDNs(fn string, startTime int64) msisdnList {
	f, err := os.Open(fn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't open '%s' fo reading: %s\n", fn, err)
		os.Exit(1)
	}
	defer f.Close()
	m, err := loadMSIDN(f, startTime-1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't initialize msidsnlist: %s\n", err)
		os.Exit(1)
	}
	if len(m) == 0 {
		fmt.Fprintf(os.Stderr, "can't work with empty msisdn list\n")
		os.Exit(1)
	}
	return m
}

func mustCreateConfig() config {

	ft := func(n, s string) int64 {
		t, err := time.Parse("2006-01-02", s)
		if err != nil {
			fmt.Fprintf(os.Stderr, "parameter %s' invalid time format '%s' (use 2019-06-20)", n, s)
			os.Exit(1)
		}
		return t.In(localZone).Unix()
	}

	startAt := flag.String("start-at", "", "start time ")
	finishAt := flag.String("finish-at", "", "finish time")
	msisdnFile := flag.String("msisdn-file", "msisdn.txt", "file with the list of caller numbers")
	rpd := flag.Int("records-per-day", 1000, "desired records per day")
	flag.Parse()

	if *rpd <= 0 {
		fmt.Fprintf(os.Stderr, "records-per-day parameter must be greter than 0")
		os.Exit(1)
	}

	c := config{
		rpd: *rpd,
	}

	for a, v := range map[string]*string{
		"start-at":    startAt,
		"finish-at":   finishAt,
		"msisdn-file": msisdnFile,
	} {
		if v == nil || *v == "" {
			fmt.Fprintf(os.Stderr, "parameter %s' is mandatory", a)
			os.Exit(1)
		}
	}

	c.startAt = ft("start-at", *startAt)
	c.finishAt = ft("finish-at", *finishAt)
	c.msisdnFile = *msisdnFile
	return c
}
