
.PHONY: run
run: binaries
	docker kill server || true
	docker rm server || true
	docker run --name server -d -p 1234:1234 --privileged server
	./cmd/client/client

.PHONY: graph
graph: graph.png
	open graph.png

graph.png: graph.gp drift.0.postprocessed.dat
	gnuplot graph.gp

drift.0.postprocessed.dat: drift.0.dat
	cat drift.0.dat | ./postprocess.py > drift.0.postprocessed.dat

.PHONY:binaries
binaries: server cmd/client/client

.PHONY: server
server:
	docker build -t server .

.PHONY: cmd/client/client
cmd/client/client:
	(cd cmd/client && go build)

.PHONY: clean
clean:
	rm -f cmd/server/server cmd/client/client *.postprocessed.dat
