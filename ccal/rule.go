package ccal

import (
	"strconv"
	"strings"
	"time"
)

type Event struct {
	*Rule
	Date time.Time
}

type Rule struct {
	Slug string
	Rule string
	Name string
}

var rules = []*Rule{
	{"s", "s", "Altjahresabend"},
	{"n", "n", "Neujahrstag"},
	{"xe", "1224", "Christnacht"},
	{"x1", "1225", "Christfest I"},
	{"x2", "1226", "Christfest II"},
	{"an", "1206", "Nikolaustag"},
	{"ep", "10106", "Epiphanias"},

	{"a1", "a", "1. Advent"},
	{"a2", "a+7", "2. Advent"},
	{"a3", "a+14", "3. Advent"},
	{"a4", "a+21", "4. Advent"},

	{"e1", "e", "1.S.n Epiphanias"},
	{"e2", "e+7", "2.S.n Epiphanias"},
	{"e3", "e+14", "3.S.n Epiphanias"},
	{"e4", "e+21", "letzter S.n Epiphanias"},

	{"pst", "p-63", "Septuagesimä"},
	{"psx", "p-56", "Sexagesimä"},
	{"pe", "p-49", "Estomihi"},
	{"pa", "p-46", "Aschermittwoch"},
	{"pi", "p-42", "Invocavit"},
	{"pr", "p-35", "Reminiszere"},
	{"po", "p-28", "Okuli"},
	{"pl", "p-21", "Lätare"},
	{"pj", "p-14", "Judika"},
	{"pp", "p-7", "Palmarum"},
	{"pg", "p-3", "Gründonnerstag"},
	{"pk", "p-2", "Karfreitag"},
	{"pn", "p-1", "Karsamstag"},
	{"p", "p", "Ostersonntag"},
	{"pt", "p+1", "Ostermontag"},
	{"p1", "p+7", "Quasimodogeniti"},
	{"p2", "p+14", "Misericordias Domini"},
	{"p3", "p+21", "Jubilate"},
	{"p4", "p+28", "Kantate"},
	{"p5", "p+35", "Rogate"},
	{"ph", "p+39", "Himmelfahrt"},
	{"p6", "p+42", "Exaudi"},
	{"pf", "p+49", "Pfingstsonntag"},
	{"pm", "p+50", "Pfingstmontag"},

	{"t", "t", "Trinitatis"},
	{"h", "h", "Erntedank"},
	{"z", "z", "Nacht zum Advent"},
	{"tj", "10624", "Johannestag"},
	{"tm", "10929", "Michaelistag"},
	{"tr", "11031", "Reformationsfest"},
	{"tm", "11111", "Martinstag"},

	{"z3", "z-20", "drittletzter S.d.Kj"},
	{"z2", "z-13", "vorletzter S.d.Kj"},
	{"zb", "z-10", "Buß- und Bettag"},
	{"z1", "z-6", "Ewigkeitssonntag"},

	{"t1", "t+7", "1.S.n Trinitatis"},
	{"t2", "t+14", "2.S.n Trinitatis"},
	{"t3", "t+21", "3.S.n Trinitatis"},
	{"t4", "t+28", "4.S.n Trinitatis"},
	{"t5", "t+35", "5.S.n Trinitatis"},
	{"t6", "t+42", "6.S.n Trinitatis"},
	{"t7", "t+49", "7.S.n Trinitatis"},
	{"t8", "t+56", "8.S.n Trinitatis"},
	{"t9", "t+63", "9.S.n Trinitatis"},
	{"t10", "t+70", "10.S.n Trinitatis"},
	{"t11", "t+77", "11.S.n Trinitatis"},
	{"t12", "t+84", "12.S.n Trinitatis"},
	{"t13", "t+91", "13.S.n Trinitatis"},
	{"t14", "t+98", "14.S.n Trinitatis"},
	{"t15", "t+105", "15.S.n Trinitatis"},
	{"t16", "t+112", "16.S.n Trinitatis"},
	{"t17", "t+119", "17.S.n Trinitatis"},
	{"t18", "t+126", "18.S.n Trinitatis"},
	{"t19", "t+133", "19.S.n Trinitatis"},
	{"t20", "t+140", "20.S.n Trinitatis"},
	{"t21", "t+147", "21.S.n Trinitatis"},
	{"t22", "t+154", "22.S.n Trinitatis"},
	{"t23", "t+161", "23.S.n Trinitatis"},
}

func (c *Ccal) Exec(rule string) time.Time {
	if rule == "" {
		return time.Time{}
	}
	fst, rst := rule[0], rule
	if fst >= '0' && fst <= '9' {
		y := 0
		if len(rule) == 5 {
			y = int(fst) - '0'
			rst = rule[1:]
		}
		if len(rst) == 4 {
			m, _ := strconv.Atoi(rst[:2])
			d, _ := strconv.Atoi(rst[2:4])
			return date(y+c.y, m, d)
		}
		return time.Time{}
	}
	idx := strings.IndexAny(rule, "-+")
	var op byte
	var arg string
	if idx > 0 {
		op = rule[idx]
		rst, arg = rule[:idx], rule[idx+1:]
	}
	t := c.Day(rst)
	if op != 0 {
		d, _ := strconv.Atoi(arg)
		switch op {
		case '+':
			t = t.AddDate(0, 0, d)
		case '-':
			t = t.AddDate(0, 0, -d)
		}
	}
	return t
}
