#!/usr/bin/env python

import sys
import random
import time

def generateHexForBytes(num_bytes, generator=None):
    values = []
    if not generator:
        generator = lambda: random.randint(0,255)
    for x in range(num_bytes):
        values.append( '%02x' % (generator(),))
    return ''.join(values)


for x in xrange(1000):
    taskid = generateHexForBytes(8)
    for y in xrange(10000):
        opid = generateHexForBytes(8)

        sys.stdout.write("""
X-Trace: %s%s
Agent: Python Test Generator
Timestamp: %f
""" % (taskid, opid, time.time()))
