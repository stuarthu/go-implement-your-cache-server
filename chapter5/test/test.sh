./cache-benchmark -type tcp -n 100000 -r 100000 -t set

./cache-benchmark -type tcp -n 100000 -r 100000 -t set -P 3

./cache-benchmark -type tcp -n 100000 -r 100000 -t set -c 50

redis-benchmark -c 50 -n 100000 -d 1000 -t set -r 100000 -P 1
