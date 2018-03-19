./cache-benchmark -type tcp -n 100000 -r 100000 -t get -P 10

./cache-benchmark -type tcp -n 100000 -r 100000 -t get -P 100

./cache-benchmark -type tcp -n 100000 -r 100000 -t get -c 50

redis-benchmark -c 50 -n 100000 -d 1000 -t get -r 100000 -P 1
