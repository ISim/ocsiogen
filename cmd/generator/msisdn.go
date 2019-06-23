package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// represents infinity validity time
var infTime = time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC).Unix()

// allocations
var allocations = []struct{ from, n int64 }{
	{311220000, 10000}, {490500000, 10000}, {511130000, 10000}, {461030000, 10000},
	{790580000, 10000}, {790040000, 20000},
}

type validity struct {
	from, to int64
}

type validities []validity

type msisdnMap map[string]validities

func loadMSIDN(in io.Reader, t int64) (msisdnMap, error) {
	nums := msisdnMap{}

	scanner := bufio.NewScanner(in)
	line := 0
	for scanner.Scan() {
		line++
		parts := strings.Split(scanner.Text(), `;`)
		if len(parts) == 0 {
			continue
		}

		v := validity{}

		if len(parts) > 1 {
			t, err := time.Parse(time.RFC3339, parts[1])
			if err != nil {
				log.Printf("dtFrom parse error line %d: %s", line, err)
				continue
			}
			v.from = t.Unix()
		}
		if len(parts) > 2 && parts[2] != "" {
			t, err := time.Parse(time.RFC3339, parts[2])
			if err != nil {
				log.Printf("dtTo parse error line %d: %s", line, err)
				continue
			}
			v.to = t.Unix()
		} else {
			v.to = infTime
		}
		if v.to <= t {
			// already dead number
			continue
		}
		if _, ok := nums[parts[0]]; !ok {
			nums[parts[0]] = validities{}
		}
		nums[parts[0]] = append(nums[parts[0]], v)
	}

	if err := scanner.Err(); err != nil {
		return nil, errors.Wrapf(err, "read error")
	}
	return nums, nil
}

func (m msisdnMap) newMsisdn(from int64) (string, int64, int64) {
	var (
		to  int64
		num string
	)

	for {
		x := rand.Intn(10)
		switch {
		case x == 9:
			to = from + 86400 + rand.Int63n(10*86400)
		case x > 6:
			to = from + 10*86400 + rand.Int63n(100*86400)
		default:
			to = infTime
		}

		if rand.Intn(100) == 99 {
			num = m.migratedNumber(from, to)
		}
		if num == "" {
			num = newNumber()
			if _, ok := m[num]; ok {
				// already exists
				continue
			}
		}
		m.appendNum(num, from, to)
		break
	}
	return num, from, to
}

func newNumber() string {
	i := rand.Intn(len(allocations))

	return fmt.Sprintf("%d", allocations[i].from+rand.Int63n(allocations[i].n))
}

// tries to find gap in validity
func (m msisdnMap) migratedNumber(from, to int64) string {
	num := ""

M:
	for n, vs := range m {
		for _, v := range vs {
			if (from >= v.from && from <= v.to) ||
				(to >= v.from && to <= v.to) ||
				(from <= v.from && to >= v.to) {
				// overlaping interval
				continue M
			}
		}
		num = n
		break
	}
	return num
}

func (m msisdnMap) appendNum(num string, from, to int64) {
	if _, ok := m[num]; !ok {
		m[num] = validities{}
	}
	m[num] = append(m[num], validity{from: from, to: to})
}

func (m msisdnMap) flush(w io.Writer) {
	for n, vs := range m {
		for _, v := range vs {
			fmt.Fprintf(w, "%s;%s;", n, time.Unix(v.from, 0).Format(time.RFC3339))
			if v.to != infTime {
				fmt.Fprintf(w, "%s", time.Unix(v.to, 0).Format(time.RFC3339))
			}
			fmt.Fprintf(w, "\n")
		}
	}
}
