set terminal png
set output 'graph.png'
set ylabel "Clock drift / milliseconds"
set xlabel "Time / seconds"
plot 'drift.0.postprocessed.dat' with lines title "Host - VM clock drift"
