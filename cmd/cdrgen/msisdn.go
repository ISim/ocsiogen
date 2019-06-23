package main

import (
	"bufio"
	"io"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type msisdn struct {
	msisdn             string
	from, to           int64
	lastCallFinishedAt int64
}

type msisdnList []*msisdn

func loadMSIDN(in io.Reader, t int64) (msisdnList, error) {
	l := msisdnList{}

	scanner := bufio.NewScanner(in)
	line := 0
	for scanner.Scan() {
		line++
		parts := strings.Split(scanner.Text(), `;`)
		if len(parts) == 0 {
			continue
		}
		m := msisdn{
			msisdn:             parts[0],
			lastCallFinishedAt: t,
		}
		if len(parts) > 1 {
			t, err := time.Parse(time.RFC3339, parts[1])
			if err != nil {
				log.Printf("dtFrom parse error line %d: %s", line, err)
				continue
			}
			m.from = t.UnixNano()
			m.lastCallFinishedAt = m.from
		}
		if len(parts) > 2 && parts[2] != "" {
			t, err := time.Parse(time.RFC3339, parts[2])
			if err != nil {
				log.Printf("dtTo parse error line %d: %s", line, err)
				continue
			}
			m.to = t.UnixNano()
		}
		if m.to != 0 && m.to <= t {
			// already dead number
			continue
		}
		l = append(l, &m)
	}

	if err := scanner.Err(); err != nil {
		return nil, errors.Wrapf(err, "read error")
	}
	return l, nil
}

func (m msisdnList) allocate(from, to int64) string {
	recs := len(m)
	idx := rand.Intn(recs)
	i := idx
	for {
		mrec := m[i]
		if mrec.from <= from &&
			(mrec.to == 0 || mrec.to >= to) &&
			mrec.lastCallFinishedAt < from {
			mrec.lastCallFinishedAt = to
			return mrec.msisdn
		}
		i++
		if i == recs {
			i = 0
		}
		if i == idx {
			// no suitable msidn found
			break
		}
	}

	return ""
}
