# ./server

curl 127.0.0.1:12345/cache/a -XPUT -daa

curl 127.0.0.1:12345/cache/a

curl 127.0.0.1:12345/status

# wait 30 seconds

curl 127.0.0.1:12345/cache/a

curl 127.0.0.1:12345/status
