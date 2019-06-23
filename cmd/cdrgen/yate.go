package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/xid"
)

type yateWriter struct {
	dir, fileTemplate string
	f                 io.WriteCloser
	hour              int64
}

func newYateWriter(dir, fileTemplate string) *yateWriter {
	return &yateWriter{
		dir:          dir,
		fileTemplate: fileTemplate,
	}
}

func (y *yateWriter) Write(c cdr) error {
	h := c.finished / nano / 3600
	if h > y.hour {
		err := y.open(time.Unix(c.finished/nano, 0))
		if err != nil {
			return err
		}
		y.hour = h
	}
	billID := xid.New().String()
	billTime := float64(c.finished-c.started) / nano
	var outStatus, inStatus, inReason, outReason, flags string
	if billTime > 0 {
		outStatus = "answered"
		inStatus = "answered"
		flags = "PS=70769,OS=11323040,PR=70758,OR=11321280,PL=0"
	} else {
		outStatus = "hangup"
		inStatus = "ringing"
		inReason = "Request terminated"
		outReason = "Cancelled"
	}
	ringTime := rand.Float64() * 30
	endTime := float64(c.finished) / nano
	_, err := fmt.Fprintf(y.f, "%.3f\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%.3f\t%.3f\t%.3f\t%s\t%s\t%s\t%s\n",
		endTime,
		"call",
		"general",
		billID,
		"sip/2222",
		"127.0.0.1/5060",
		c.msisdn,
		c.bnum,
		billTime,
		ringTime,
		billTime+ringTime,
		"incoming",
		inStatus,
		inReason,
		flags,
	)
	if err != nil {
		return err
	}
	endTime += 0.005 - rand.Float64()/1000
	_, err = fmt.Fprintf(y.f, "%.3f\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%.3f\t%.3f\t%.3f\t%s\t%s\t%s\t%s\n",
		endTime,
		"call",
		"artf_trunk",
		billID,
		"sip/2222",
		"127.0.0.1/5060",
		c.msisdn,
		c.bnum,
		billTime,
		ringTime,
		billTime+ringTime+0.005-rand.Float64()/1000,
		"outgoing",
		outStatus,
		outReason,
		flags,
	)

	return err
}

func (y *yateWriter) Close() {
	if y.f != nil {
		y.f.Close()
		y.f = nil
	}
}

func (y *yateWriter) open(t time.Time) error {
	var err error
	if y.f != nil {
		err = y.f.Close()
		y.f = nil
	}
	if err != nil {
		return err
	}
	p := filepath.Join(y.dir, fmt.Sprintf("%s%04d%02d%02d%02d%02d",
		y.fileTemplate,
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute()))
	f, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return errors.Wrapf(err, "can't open file '%s'", p)
	}
	y.f = f
	return nil
}
