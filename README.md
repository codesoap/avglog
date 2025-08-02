avglog takes a file with lines of the format `yyyy-mm-ddThh:mm <number>`.
It averages values over a given duration and adds these averages as a
third column.

# Installation
```bash
go install github.com/codesoap/avglog@latest
```

# Usage
```console
$ avglog -h
Usage: avglog [-d <window-size>] [<input-file>]
        If <input-file> is not given, standard input will be read.
Options:
  -d duration
        averaging window size (default 240h0m0s)

$ # Average with the default window of ten days:
$ head bodyweight.log
2025-01-01T16:14 56.7
2025-01-02T17:01 56.9
2025-01-04T16:47 56.6
2025-01-05T18:30 57.0
$ avglog bodyweight.log > bodyweight_with_avg.log
$ head bodyweight_with_avg.log
2025-01-01T16:14 56.7 56.8
2025-01-02T17:01 56.9 56.81999999999999
2025-01-04T16:47 56.6 56.849999999999994
2025-01-05T18:30 57 56.925

$ # Average with a window of 30 days:
$ avglog -d 720h watertemp.log > watertemp_with_avg.log
```

# Printing with gnuplot
You can use a gnuplot script similar to this one to plot the data with
its average into a PNG file:

```
set grid
set xdata time
set timefmt "%Y-%m-%dT%H:%M"
set format x "%Y-%m"
set ytics 1

set terminal pngcairo size 1800,600 enhanced font "Helvetica,20"
set output 'bodyweight.png'
plot 'bodyweight_with_avg.log' using 1:2 w points pointtype 7 lw 2, \
     'bodyweight_with_avg.log' using 1:3 w lines lw 4 lc 'orange'
```
