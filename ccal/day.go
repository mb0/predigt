package ccal

import (
	"sort"
	"time"
)

type Ccal struct {
	y       int
	a, p, z time.Time
}

func NewCcal(year int) *Ccal {
	return &Ccal{y: year,
		a: Advent(year),
		p: Computus(year + 1),
		z: Advent(year+1).AddDate(0, 0, -1),
	}
}

func (c *Ccal) All() []Event {
	res := make([]Event, 0, len(rules))
	for _, r := range rules {
		d := c.Exec(r.Rule)
		idx := sort.Search(len(res), func(i int) bool {
			return !d.After(res[i].Date)
		})
		if idx < len(res) {
			if res[idx].Date.Equal(d) {
				continue
			} else {
				res = append(res[:idx+1], res[idx:]...)
			}
			res[idx] = Event{Rule: r, Date: d}
		} else {
			res = append(res, Event{Rule: r, Date: d})
		}
	}
	return res
}

func (c *Ccal) Day(n string) time.Time {
	switch n {
	case "a":
		return c.a
	case "s":
		return date(c.y+1, 1, 1).AddDate(0, 0, -1)
	case "n":
		return date(c.y+1, 1, 1)
	case "e":
		return Epiphanias1(c.y + 1)
	case "p":
		return c.p
	case "t":
		return c.p.AddDate(0, 0, 56)
	case "z":
		return c.z
	case "h":
		d := date(c.y+1, 10, 1).AddDate(0, 0, -1)
		w := int(d.Weekday())
		return d.AddDate(0, 0, 7-w)
	}
	return time.Time{}
}

func date(y, m, d int) time.Time {
	return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.Local)
}
func split(t time.Time) (y, m, d int) {
	return t.Year(), int(t.Month()), t.Day()
}

// Computus returns the easter day for the given year using
// https://en.wikipedia.org/wiki/Date_of_Easter#Anonymous_Gregorian_algorithm
func Computus(year int) time.Time {
	if year < 1583 || year > 9999 {
		return time.Time{}
	}
	a := year % 19
	b := year / 100
	c := year % 100
	d := b / 4
	e := b % 4
	f := (b + 8) / 25
	g := (b - f + 1) / 3
	h := (19*a + b - d - g + 15) % 30
	i := c / 4
	k := c % 4
	l := (32 + 2*e + 2*i - h - k) % 7
	m := (a + 11*h + 22*l) / 451
	n := (h + l - 7*m + 114) / 31
	o := (h + l - 7*m + 114) % 31
	return date(year, n, o+1)
}

func Advent(year int) time.Time {
	t := date(year, 12, 24)
	w := int(t.Weekday())
	return t.AddDate(0, 0, -w-3*7)
}
func Epiphanias1(year int) time.Time {
	t := date(year, 1, 6)
	w := int(t.Weekday())
	return t.AddDate(0, 0, 7-w)
}
