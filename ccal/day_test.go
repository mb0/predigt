package ccal

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestCcalExec(t *testing.T) {
	tests := []struct {
		year       int
		rule, want string
	}{
		{2023, "1225", "2023-12-25"},
		{2023, "p", "2024-03-31"},
		{2023, "p-3", "2024-03-28"},
		{2023, "t+119", "2024-09-22"},
		{2023, "11111", "2024-11-11"},
	}
	for _, test := range tests {
		cc := NewCcal(test.year)
		d := cc.Exec(test.rule)
		got := d.Format("2006-01-02")
		if got != test.want {
			t.Errorf("%s want %s got %s", test.rule, test.want, got)
		}
	}
}

func TestCcalAll(t *testing.T) {
	testAll(t, 2022, year22Raw)
	testAll(t, 2023, year23Raw)
}

func testAll(t *testing.T, year int, raw string) {
	cc := NewCcal(year)
	all := cc.All()
	if len(all) == 0 {
		t.Errorf("empty events for year %d", year)
		return
	}
	lines := strings.Split(strings.TrimSpace(raw), "\n")

	var j int
	for _, want := range lines {
		if j >= len(all) {
			t.Errorf("year %d want rest %s…", year, want)
			break
		}
		ev := all[j]
		got := fmt.Sprintf("%s %s", ev.Date.Format("20060102"), ev.Name)
		if got != want {
			t.Errorf("year %d for %s want %s got %s", year, ev.Slug, want, got)
		}
		j++
	}
}

func TestDays(t *testing.T) {
	testDay(t, "advent", adventRaw, Advent)
	testDay(t, "eastern", easternRaw, Computus)
}

func testDay(t *testing.T, name, raw string, f func(int) time.Time) {
	lines := strings.Split(raw, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, " ")
		for _, part := range parts {
			want, err := time.ParseInLocation("02.01.2006", part, time.Local)
			if err != nil {
				t.Errorf("%s for %s err %v", name, part, err)
			}
			got := f(want.Year())
			if got != want {
				t.Errorf("%s for %s want %s got %s", name, part, want, got)
			}
		}
	}
}

var adventRaw = `29.11.2020 28.11.2021 27.11.2022 03.12.2023 01.12.2024 30.11.2025 29.11.2026`

var easternRaw = `
23.03.1913 24.03.1940 25.03.1951 26.03.1967 27.03.1910 28.03.1937 29.03.1959 30.03.1902 31.03.1907
23.03.2008 25.03.2035 26.03.1978 27.03.1921 28.03.1948 29.03.1964 30.03.1975 31.03.1918
25.03.2046 26.03.1989 27.03.1932 28.03.2027 29.03.1970 30.03.1986 31.03.1929
26.03.2062 27.03.2005 28.03.2032 29.03.2043 30.03.1997 31.03.1991
26.03.2073 27.03.2016 28.03.2100 29.03.2054 30.03.2059 31.03.2002
26.03.2084 29.03.2065 30.03.2070 31.03.2013
30.03.2081 31.03.2024
30.03.2092 31.03.2086
31.03.2097
01.04.1923 02.04.1961 03.04.1904 04.04.1915 05.04.1931 06.04.1947 07.04.1901 08.04.1917 09.04.1939
01.04.1934 02.04.1972 03.04.1983 04.04.1920 05.04.1942 06.04.1958 07.04.1912 08.04.1928 09.04.1944
01.04.1945 02.04.2051 03.04.1988 04.04.1926 05.04.1953 06.04.1969 07.04.1985 08.04.2007 09.04.1950
01.04.1956 02.04.2056 03.04.1994 04.04.1999 05.04.2015 06.04.1980 07.04.1996 08.04.2012 09.04.2023
01.04.2018 03.04.2067 04.04.2010 05.04.2026 06.04.2042 07.04.2075 08.04.2091 09.04.2034
01.04.2029 03.04.2078 04.04.2021 05.04.2037 06.04.2053 07.04.2080 09.04.2045
01.04.2040 03.04.2089 04.04.2083 05.04.2048 06.04.2064
04.04.2094
10.04.1955 11.04.1909 12.04.1903 13.04.1941 14.04.1963 15.04.1900 16.04.1911 17.04.1927 18.04.1954
10.04.1966 11.04.1971 12.04.1914 13.04.1952 14.04.1968 15.04.1906 16.04.1922 17.04.1938 18.04.1965
10.04.1977 11.04.1982 12.04.1925 13.04.2031 14.04.1974 15.04.1979 16.04.1933 17.04.1949 18.04.1976
10.04.2039 11.04.1993 12.04.1936 13.04.2036 14.04.2047 15.04.1990 16.04.1995 17.04.1960 18.04.2049
10.04.2050 11.04.2004 12.04.1998 14.04.2058 15.04.2001 16.04.2006 17.04.2022 18.04.2055
10.04.2061 11.04.2066 12.04.2009 14.04.2069 15.04.2063 16.04.2017 17.04.2033 18.04.2060
10.04.2072 11.04.2077 12.04.2020 15.04.2074 16.04.2028 17.04.2044
11.04.2088 12.04.2093 15.04.2085 16.04.2090
12.04.2099 15.04.2096
19.04.1908 20.04.1919 21.04.1935 22.04.1962 23.04.1905 24.04.2011 25.04.1943
19.04.1981 20.04.1924 21.04.1946 22.04.1973 23.04.1916 24.04.2095 25.04.2038
19.04.1987 20.04.1930 21.04.1957 22.04.1973 23.04.2000
19.04.1992 20.04.2003 21.04.2019 22.04.2057 23.04.2079
19.04.2071 20.04.2014 21.04.2030 22.04.2068
19.04.2076 20.04.2025 21.04.2041
19.04.2082 20.04.2087 21.04.2052
`
var year22Raw = `
20221127 1. Advent
20221204 2. Advent
20221206 Nikolaustag
20221211 3. Advent
20221218 4. Advent
20221224 Christnacht
20221225 Christfest I
20221226 Christfest II
20221231 Altjahresabend
20230101 Neujahrstag
20230106 Epiphanias
20230108 1.S.n Epiphanias
20230115 2.S.n Epiphanias
20230122 3.S.n Epiphanias
20230129 letzter S.n Epiphanias
20230205 Septuagesimä
20230212 Sexagesimä
20230219 Estomihi
20230222 Aschermittwoch
20230226 Invocavit
20230305 Reminiszere
20230312 Okuli
20230319 Lätare
20230326 Judika
20230402 Palmarum
20230406 Gründonnerstag
20230407 Karfreitag
20230408 Karsamstag
20230409 Ostersonntag
20230410 Ostermontag
20230416 Quasimodogeniti
20230423 Misericordias Domini
20230430 Jubilate
20230507 Kantate
20230514 Rogate
20230518 Himmelfahrt
20230521 Exaudi
20230528 Pfingstsonntag
20230529 Pfingstmontag
20230604 Trinitatis
20230611 1.S.n Trinitatis
20230618 2.S.n Trinitatis
20230624 Johannestag
20230625 3.S.n Trinitatis
20230702 4.S.n Trinitatis
20230709 5.S.n Trinitatis
20230716 6.S.n Trinitatis
20230723 7.S.n Trinitatis
20230730 8.S.n Trinitatis
20230806 9.S.n Trinitatis
20230813 10.S.n Trinitatis
20230820 11.S.n Trinitatis
20230827 12.S.n Trinitatis
20230903 13.S.n Trinitatis
20230910 14.S.n Trinitatis
20230917 15.S.n Trinitatis
20230924 16.S.n Trinitatis
20230929 Michaelistag
20231001 Erntedank
20231008 18.S.n Trinitatis
20231015 19.S.n Trinitatis
20231022 20.S.n Trinitatis
20231029 21.S.n Trinitatis
20231031 Reformationsfest
20231105 22.S.n Trinitatis
20231111 Martinstag
20231112 drittletzter S.d.Kj
20231119 vorletzter S.d.Kj
20231122 Buß- und Bettag
20231126 Ewigkeitssonntag
`
var year23Raw = `
20231203 1. Advent
20231206 Nikolaustag
20231210 2. Advent
20231217 3. Advent
20231224 Christnacht
20231225 Christfest I
20231226 Christfest II
20231231 Altjahresabend
20240101 Neujahrstag
20240106 Epiphanias
20240107 1.S.n Epiphanias
20240114 2.S.n Epiphanias
20240121 3.S.n Epiphanias
20240128 letzter S.n Epiphanias
20240204 Sexagesimä
20240211 Estomihi
20240214 Aschermittwoch
20240218 Invocavit
20240225 Reminiszere
20240303 Okuli
20240310 Lätare
20240317 Judika
20240324 Palmarum
20240328 Gründonnerstag
20240329 Karfreitag
20240330 Karsamstag
20240331 Ostersonntag
20240401 Ostermontag
20240407 Quasimodogeniti
20240414 Misericordias Domini
20240421 Jubilate
20240428 Kantate
20240505 Rogate
20240509 Himmelfahrt
20240512 Exaudi
20240519 Pfingstsonntag
20240520 Pfingstmontag
20240526 Trinitatis
20240602 1.S.n Trinitatis
20240609 2.S.n Trinitatis
20240616 3.S.n Trinitatis
20240623 4.S.n Trinitatis
20240624 Johannestag
20240630 5.S.n Trinitatis
20240707 6.S.n Trinitatis
20240714 7.S.n Trinitatis
20240721 8.S.n Trinitatis
20240728 9.S.n Trinitatis
20240804 10.S.n Trinitatis
20240811 11.S.n Trinitatis
20240818 12.S.n Trinitatis
20240825 13.S.n Trinitatis
20240901 14.S.n Trinitatis
20240908 15.S.n Trinitatis
20240915 16.S.n Trinitatis
20240922 17.S.n Trinitatis
20240929 Michaelistag
20241006 Erntedank
20241013 20.S.n Trinitatis
20241020 21.S.n Trinitatis
20241027 22.S.n Trinitatis
20241031 Reformationsfest
20241103 23.S.n Trinitatis
20241110 drittletzter S.d.Kj
20241111 Martinstag
20241117 vorletzter S.d.Kj
20241120 Buß- und Bettag
20241124 Ewigkeitssonntag
`
