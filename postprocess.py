#!/usr/bin/env python

import sys

# title line
line = sys.stdin.readline().strip()
print line

# first data point
line = sys.stdin.readline().strip()
bits = line.split()
first = int(bits[1])
print bits[0], "0"

while True:
	line = sys.stdin.readline().strip()
	if line == "":
		sys.exit(0)
	bits = line.split()
	bits[1] = str((int(bits[1]) - first) / 1000)
        if len(bits) > 2:
            bits[2] = str(int(bits[2]) / 1000)
	print " ".join(bits)

