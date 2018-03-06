BIN = test_basic test_batch_write test_read_throughput ingest_data test_write_throughput test_read test_property test_ttl
CXXFLAGS = -I /home/stuart/workspace/go-implement-your-cache-server/rocksdb/include/ -std=gnu++14 -L ~/workspace/go-implement-your-cache-server/rocksdb -lrocksdb -lpthread -lz -lsnappy -lboost_program_options

ifeq ($(MAKECMDGOALS), debug)
    CXXFLAGS += -g -ggdb
else
    CXXFLAGS += -O3
endif

all: release

debug: $(BIN)

release: $(BIN)

test_basic: test_basic.cpp common.h
	g++ -o $@ $< $(CXXFLAGS)

test_ttl: test_ttl.cpp common.h
	g++ -o $@ $< $(CXXFLAGS)

test_property: test_property.cpp common.h
	g++ -o $@ $< $(CXXFLAGS)

test_batch_write: test_batch_write.cpp common.h
	g++ -o $@ $< $(CXXFLAGS)

test_read: test_read.cpp common.h
	g++ -o $@ $< $(CXXFLAGS)

test_read_throughput: test_read_throughput.cpp common.h
	g++ -o $@ $< $(CXXFLAGS)

test_write_throughput: test_write_throughput.cpp common.h
	g++ -o $@ $< $(CXXFLAGS)

ingest_data: ingest_data.cpp common.h
	g++ -o $@ $< $(CXXFLAGS)

clean:
	rm -f $(BIN)
