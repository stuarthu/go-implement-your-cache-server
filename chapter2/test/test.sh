./client -c set -k testkey -v testvalue

./client -c get -k testkey

curl 127.0.0.1:12345/status

./client -c del -k testkey

curl 127.0.0.1:12345/status

./cache-benchmark -type tcp -n 100000 -r 100000 -t set

./cache-benchmark -type tcp -n 100000 -r 100000 -t get
