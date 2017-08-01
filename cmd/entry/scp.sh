#!/bin/sh
tar -czf entry.tar.gz entry
scp -P 28315 entry.tar.gz root@duomila.club:/root/ 
