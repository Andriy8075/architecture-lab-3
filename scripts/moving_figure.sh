#!/bin/bash
curl -X POST http://localhost:17000 -d "figure 0.5 0.5"
curl -X POST http://localhost:17000 -d "update"
sleepTime = 0.1

for i in {1..3}; do
    curl -X POST http://localhost:17000 -d "move 0.3 0"
    curl -X POST http://localhost:17000 -d "update"
    sleep sleepTime
    curl -X POST http://localhost:17000 -d "move 0 0.3"
    curl -X POST http://localhost:17000 -d "update"
    sleep sleepTime
    curl -X POST http://localhost:17000 -d "move -0.3 0"
    curl -X POST http://localhost:17000 -d "update"
    sleep sleepTime
    curl -X POST http://localhost:17000 -d "move 0 -0.3"
    curl -X   POST http://localhost:17000 -d "update"
    sleep sleepTime
done