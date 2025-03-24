#!/bin/bash
curl -X POST http://localhost:17000 -d "white"
curl -X POST http://localhost:17000 -d "figure 0.5 0.5"
curl -X POST http://localhost:17000 -d "update"

for i in {1..10}; do
    x=$(echo "0.1 * $i" | bc)
    y=$(echo "0.1 * $i" | bc)
    curl -X POST http://localhost:17000 -d "move $x $y"
    curl -X POST http://localhost:17000 -d "update"
    sleep 1
done