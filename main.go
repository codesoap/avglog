package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	window time.Duration
	in     io.Reader
)

type measurement struct {
	t   time.Time
	num float64
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-d <window-size>] [<input-file>]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\tIf <input-file> is not given, standard input will be read.\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.DurationVar(&window, "d", 240*time.Hour, "averaging window size")
	flag.Parse()
	switch len(flag.Args()) {
	case 0:
		in = os.Stdin
	case 1:
		f, err := os.Open(flag.Arg(0))
		if err != nil {
			format := "Could not open file '%s': %s\n"
			fmt.Fprintf(os.Stderr, format, flag.Arg(0), err)
			os.Exit(1)
		}
		in = f
	default:
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	printWithAvg(measurements())
}

func measurements() []measurement {
	var msrs []measurement
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) != 2 {
			format := "Did not find exactly two columns in line '%s'.\n"
			fmt.Fprintf(os.Stderr, format, scanner.Text())
			os.Exit(1)
		}
		t, err := time.Parse("2006-01-02T15:04", fields[0])
		if err != nil {
			format := "Invalid timestamp in line '%s': %s\n"
			fmt.Fprintf(os.Stderr, format, scanner.Text(), err)
			os.Exit(1)
		}
		num, err := strconv.ParseFloat(fields[1], 64)
		if err != nil {
			format := "Invalid number in line '%s': %s\n"
			fmt.Fprintf(os.Stderr, format, scanner.Text(), err)
			os.Exit(1)
		}
		msrs = append(msrs, measurement{t: t, num: num})
	}
	if scanner.Err() != nil {
		fmt.Fprintln(os.Stderr, "Could read a line:", scanner.Err())
		os.Exit(1)
	}
	return msrs
}

func printWithAvg(msrs []measurement) {
	sortedMsrs := sortMsrs(msrs)
	nums := make([]float64, 10)
	for _, msr := range msrs {
		nums = nums[:0]
		since := msr.t.Add(-window / 2)
		until := msr.t.Add(window / 2)
		for _, m := range sortedMsrs {
			if m.t.After(until) {
				break
			}
			if m.t.After(since) {
				nums = append(nums, m.num)
			}
		}
		var sum float64
		for _, num := range nums {
			sum += num
		}
		avg := sum / float64(len(nums))
		fmt.Println(msr.t.Format("2006-01-02T15:04"), msr.num, avg)
	}
}

func sortMsrs(msrs []measurement) []measurement {
	sortedMsrs := make([]measurement, len(msrs))
	for i := range msrs {
		sortedMsrs[i] = msrs[i]
	}
	sort.Slice(sortedMsrs, func(i, j int) bool {
		return sortedMsrs[i].t.Before(sortedMsrs[j].t)
	})
	return sortedMsrs
}
