# ./server -node 1.1.1.1

# ./server -node 1.1.1.2 -cluster 1.1.1.1

./client -h 1.1.1.1 -c set -k keya -v a

./client -h 1.1.1.1 -c set -k keyb -v b

./client -h 1.1.1.1 -c set -k keyc -v c

./client -h 1.1.1.1 -c set -k keyd -v d

./client -h 1.1.1.1 -c set -k keye -v e

# ./server -node 1.1.1.3 -cluster 1.1.1.2

./client -h 1.1.1.1 -c set -k keya -v a

./client -h 1.1.1.1 -c set -k keyb -v b

./client -h 1.1.1.1 -c set -k keyc -v c

./client -h 1.1.1.1 -c set -k keyd -v d

./client -h 1.1.1.1 -c set -k keye -v e

# stop 1.1.1.1

./client -h 1.1.1.2 -c set -k keya -v a

./client -h 1.1.1.2 -c set -k keyb -v b

./client -h 1.1.1.2 -c set -k keyc -v c

./client -h 1.1.1.2 -c set -k keyd -v d

./client -h 1.1.1.2 -c set -k keye -v e
