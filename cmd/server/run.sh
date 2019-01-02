#!/bin/sh
docker run -it -p 1234:1234 -v `pwd`:/src alpine /src/server
