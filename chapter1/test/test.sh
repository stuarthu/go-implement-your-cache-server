curl 127.0.0.1:12345/status

curl -v 127.0.0.1:12345/cache/testkey -XPUT -dtestvalue

curl 127.0.0.1:12345/cache/testkey

curl 127.0.0.1:12345/status

curl 127.0.0.1:12345/cache/testkey -XDELETE

curl 127.0.0.1:12345/status

./cache-benchmark -type http -n 100000 -r 100000 -t set

./cache-benchmark -type http -n 100000 -r 100000 -t get

redis-benchmark -c 1 -n 100000 -d 1000 -t set,get -r 100000
