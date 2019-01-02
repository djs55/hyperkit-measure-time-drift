# hyperkit-measure-time-drift
Simple tool to measure the time drift between a VM and the host

While running Docker Desktop for Mac, start gathering data:

```
$ make run
docker kill server || true
server
docker rm server || true
server
docker run --name server -d -p 1234:1234 --privileged server
3708b07b82be7c6472ec00937c179a3952f0be5931a357311339c643cde8e91e
./cmd/client/client
2019/01/02 11:24:59 Created drift.0.dat
```

This will measure drift once per second.

In another terminal, draw a graph:

```
$ make graph
cat drift.0.dat | ./postprocess.py > drift.0.postprocessed.dat
gnuplot graph.gp
open graph.png
```

![sawtooth graph showing drift](https://djs55.github.io/hyperkit-measure-time-drift/screenshot.png)

Note there is no attempt to measure and subtract the round-trip time to the VM.
The graphs therefore have an arbitrary offset added and should be used to understand
the change in the drift over time, rather than the exact drift at any specific point
in time.
