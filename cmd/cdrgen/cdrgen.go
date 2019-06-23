package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
)

const nano = 1000000000

var dayDistPct = []float64{
	1,
	1,
	1,
	1,
	1,
	2,
	5,
	5,
	5,
	10,
	8,
	7,
	8,
	7,
	9,
	6,
	5,
	5,
	3,
	3,
	2,
	2,
	2,
	1,
}

type cdr struct {
	started  int64
	finished int64
	msisdn   string
	bnum     string
}

type hourCdrs []cdr

func (hc hourCdrs) Len() int {
	return len(hc)
}

func (hc hourCdrs) Less(i, j int) bool {
	return hc[i].finished < hc[j].finished
}

func (hc hourCdrs) Swap(i, j int) {
	hc[i], hc[j] = hc[j], hc[i]
}

func day(t0 int64, recs int, msisdns msisdnList) []hourCdrs {
	cdrs := make([]hourCdrs, 24)
	fract := 0.0
	for h := int64(0); h < 24; h++ {
		fract += dayDistPct[h] / 100 * float64(recs)
		cnt := int(fract)
		cdrs[h] = hour(t0+h*3600, cnt, msisdns)
		if cnt > 1 {
			fract = fract - float64(cnt)
		}
	}
	return cdrs
}

func hour(t0 int64, recs int, msisdns msisdnList) hourCdrs {
	cdrs := hourCdrs{}

	for recs > 0 {
		recs--
		c := cdr{
			finished: t0*nano + rand.Int63n(3600*nano),
		}
		c.started = c.finished - callLen()
		c.msisdn = msisdns.allocate(c.started, c.finished)
		if c.msisdn == "" {
			// all numbers are calling right now
			continue
		}
		c.bnum = bnum()
		cdrs = append(cdrs, c)
	}
	sort.Sort(cdrs)
	return cdrs
}

func callLen() int64 {
	tn := rand.Intn(1000)
	switch {
	case tn > 990:
		return 600*nano + rand.Int63n(3*3000*nano) // extra looong call
	case tn > 610:
		return 180*nano + rand.Int63n(240*nano)
	}
	l := rand.NormFloat64()*100*nano + 120*nano
	if l < 0 {
		return 0
	}
	return int64(l)
}

func bnum() string {
	tmp := rand.Intn(1000)
	switch {
	case tmp > 950:
		return extras[rand.Intn(len(extras))]
	case tmp > 700:
		prefix := intPrefixes[rand.Intn(len(intPrefixes))]
		return fmt.Sprintf("+%s%0d", prefix, rand.Intn(int(math.Pow10(12-len(prefix)))))
	}
	prefix := nationalPrefixes[rand.Intn(len(nationalPrefixes))]
	return fmt.Sprintf("%s%0d", prefix, rand.Intn(int(math.Pow10(9-len(prefix)))))

}
