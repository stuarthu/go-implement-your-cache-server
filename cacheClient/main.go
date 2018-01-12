package main

import (
	"./client"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
)

var keyPrefix = strings.Repeat("a", 99)

type statistic struct {
	count int
	time  time.Duration
}

type result struct {
	getCount   int
	missCount  int
	setCount   int
	statBucket map[string]statistic
}

func get(client client.Client, key string) (string, time.Duration) {
	start := time.Now()
	value := client.Get(key)
	return value, time.Now().Sub(start)
}

func set(client client.Client, key, value string) time.Duration {
	start := time.Now()
	client.Set(key, value)
	return time.Now().Sub(start)
}

func addStatistic(bucket map[string]statistic, key string, stat statistic) {
	b := bucket[key]
	b.count += stat.count
	b.time += stat.time
	bucket[key] = b
}

func updateStatistic(bucket map[string]statistic, d time.Duration) {
	if d < time.Millisecond {
		addStatistic(bucket, "<1ms", statistic{1, d})
	}
	if d < 10*time.Millisecond {
		addStatistic(bucket, "<10ms", statistic{1, d})
	}
	addStatistic(bucket, "total", statistic{1, d})
}

func operate(typ, server, operation string, id, count, valueSize, total int, random bool, c chan *result) {
	client := client.NewCacheClient(typ, server)
	valuePrefix := strings.Repeat("a", valueSize)
	getCount := 0
	missCount := 0
	setCount := 0
	bucket := make(map[string]statistic)
	for i := 0; i < count; i++ {
		var tmp int
		if random {
			tmp = rand.Intn(total)
		} else {
			tmp = id*count + i
		}
		key := fmt.Sprintf("%s%d", keyPrefix, tmp)
		value := fmt.Sprintf("%s%d", valuePrefix, tmp)
		var val string
		var d time.Duration
		if operation != "set" {
			val, d = get(client, key)
			if val != "" && val != value {
				log.Println(key)
				panic(val)
			}
			if val == value {
				getCount++
			} else {
				missCount++
			}
			updateStatistic(bucket, d)
		}
		if operation == "set" || (operation == "mixed" && val == "") {
			d := set(client, key, value)
			setCount++
			updateStatistic(bucket, d)
		}
	}
	c <- &result{getCount, missCount, setCount, bucket}
}

func main() {
	var typ, server, operation string
	var total, valueSize, threads int
	var random bool
	flag.StringVar(&typ, "type", "http", "cache server type")
	flag.StringVar(&server, "server", "localhost:12345", "cache server address")
	flag.IntVar(&total, "total", 10000, "total API calls")
	flag.IntVar(&valueSize, "size", 1000, "value size")
	flag.IntVar(&threads, "threads", 50, "threads")
	flag.StringVar(&operation, "operation", "set", "operation could be get/set/mixed")
	flag.BoolVar(&random, "random", false, "iterate keys in order or random")
	flag.Parse()
	fmt.Println("type is", typ)
	fmt.Println("server is", server)
	fmt.Println("total", total, "keys")
	fmt.Println("value size is", valueSize)
	fmt.Println("we have", threads, "threads")
	fmt.Println("operation is", operation)
	fmt.Println("random is", random)
	if operation != "get" && operation != "set" && operation != "mixed" {
		panic("invalid operation")
	}

	rand.Seed(time.Now().UnixNano())
	c := make(chan *result, threads)
	getCount := 0
	missCount := 0
	setCount := 0
	bucket := make(map[string]statistic)
	start := time.Now()
	for i := 0; i < threads; i++ {
		go operate(typ, server, operation, i, total/threads, valueSize, total, random, c)
	}
	for i := 0; i < threads; i++ {
		r := <-c
		getCount += r.getCount
		missCount += r.missCount
		setCount += r.setCount
		for k, v := range r.statBucket {
			addStatistic(bucket, k, v)
		}
	}
	d := time.Now().Sub(start)
	totalCount := getCount + missCount + setCount
	fmt.Printf("%d records get\n", getCount)
	fmt.Printf("%d records miss\n", missCount)
	fmt.Printf("%d records set\n", setCount)
	fmt.Printf("%d usec total\n", d/time.Microsecond)
	fmt.Printf("%d usec average for all clients\n", int64(d/time.Microsecond)/int64(totalCount))
	for k, v := range bucket {
		fmt.Printf("%d%% requests %s, %d usec average for each client\n", v.count*100/totalCount,
			k, int64(v.time/time.Microsecond)/int64(v.count))
	}
	fmt.Printf("throughput is %f MB/s\n", float64((getCount+setCount)*valueSize)/1e6/d.Seconds())
	fmt.Printf("rps is %f\n", float64(totalCount)/float64(d.Seconds()))
}
