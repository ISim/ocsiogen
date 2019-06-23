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
	startAt       int64
	finishAt      int64
	msisdnFile    string
	subscriptions int
}

func main() {
	localZone, _ = time.LoadLocation("Europe/Prague")
	cfg := mustCreateConfig()
	msisdns := mustLoadMSISDNs(cfg.msisdnFile, cfg.startAt)
	rand.Seed(time.Now().UnixNano())

	days := (cfg.finishAt-cfg.startAt)/86400 + 1
	perDay := cfg.subscriptions / int(days)
	dt := time.Unix(cfg.startAt, 0)
	subscribers := 0
	for d := days; d > 0; d-- {
		for i := 0; i < perDay; i++ {
			newSubscriber(msisdns, dt)
			subscribers++
		}

		dt = dt.Add(24 * time.Hour)
	}
	for ; subscribers < cfg.subscriptions; subscribers++ {
		newSubscriber(msisdns, dt)
	}

	writeResult(msisdns, cfg.msisdnFile)
}

func newSubscriber(msisdns msisdnMap, t time.Time) {
	dnum, validFrom, validTo := msisdns.newMsisdn(t.Unix() + rand.Int63n(86400))
	_, _, _ = dnum, validFrom, validTo
}

func writeResult(msisdns msisdnMap, fn string) {
	backup := fn + ".orig"
	os.Remove(backup)
	os.Rename(fn, backup)
	f, err := os.OpenFile(fn, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't open '%s' for writing: %s\n", fn, err)
		os.Exit(1)
	}
	defer f.Close()
	msisdns.flush(f)
	if err := f.Sync(); err != nil {
		fmt.Fprintf(os.Stderr, "can't write data into '%s': %s\n", fn, err)
		os.Exit(1)
	}
}

func mustLoadMSISDNs(fn string, startTime int64) msisdnMap {
	f, err := os.Open(fn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't open '%s' for reading: %s\n", fn, err)
		os.Exit(1)
	}
	defer f.Close()
	m, err := loadMSIDN(f, startTime-1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't initialize msidsnlist: %s\n", err)
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
	sbs := flag.Int("subscriptions", 10, "desired subscriptions number")
	flag.Parse()

	if *sbs <= 0 {
		fmt.Fprintf(os.Stderr, "records-per-day parameter must be greter than 0")
		os.Exit(1)
	}

	c := config{
		subscriptions: *sbs,
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
