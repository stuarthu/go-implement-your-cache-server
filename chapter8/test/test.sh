# ./server -node 1.1.1.1

./cache-benchmark -type tcp -n 10000 -d 1 -h 1.1.1.1

curl 1.1.1.1:12345/status

# ./server -node 1.1.1.2 -cluster 1.1.1.1

curl 1.1.1.2:12345/status

curl 1.1.1.1:12345/rebalance -XPOST

curl 1.1.1.1:12345/status

curl 1.1.1.2:12345/status

# ./server -node 1.1.1.3 -cluster 1.1.1.2

curl 1.1.1.1:12345/rebalance -XPOST

curl 1.1.1.2:12345/rebalance -XPOST

curl 1.1.1.1:12345/status

curl 1.1.1.2:12345/status

curl 1.1.1.3:12345/status
