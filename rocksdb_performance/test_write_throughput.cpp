#include "common.h"

#include <sys/time.h>

namespace {
struct threadOption {
    int id;
    int num;
    int valueSize;
    int total;
};

struct result {
    int num;
    long usec;
};

void* func(void* arg)
{
    threadOption* o = (threadOption*)arg;
    DB* db = opendb(false, "/mnt/rocksdb_" + std::to_string(o->id));
    int num = o->num;
    std::string value = std::string(o->valueSize, 'a');
    struct timeval tv;
    gettimeofday(&tv, 0);
    long start = tv.tv_sec * 1000000 + tv.tv_usec;
    for (int i = 0; i < num; i++) {
        int r = rand() % o->total;
        std::string key = keyPrefix + std::to_string(r);
        db->Put(WriteOptions(), key, value);
    }
    gettimeofday(&tv, 0);
    long end = tv.tv_sec * 1000000 + tv.tv_usec;
    result* r = new result;
    r->num = num;
    r->usec = end - start;
    delete db;
    return r;
}
}
int main(int argc, char** argv)
{
    int threads;
    options_description threadsOption("thread option");
    threadsOption.add_options()("thread,r", value<int>(&threads)->default_value(1), "thread number");

    variables_map vm = parse(argc, argv, &threadsOption);
    int total = vm["total"].as<int>();
    int valueSize = vm["size"].as<int>();
    std::cout << "thread number is " << threads << std::endl;

    pthread_t get_thread[threads];
    int id = 0;
    struct timeval tv;
    gettimeofday(&tv, 0);
    long start = tv.tv_sec * 1000000 + tv.tv_usec;
    for (int i = 0; i < threads; i++) {
        threadOption* o = new threadOption;
        o->id = id++;
        o->num = total / threads;
        o->valueSize = valueSize;
        o->total = total;
        pthread_create(&get_thread[i], 0, func, o);
    }

    for (int i = 0; i < threads; i++) {
        void* res;
        pthread_join(get_thread[i], &res);
        result* r = (result*)res;
        std::cout << r->num << " records get in " << r->usec << " usec, "
                  << double(r->usec) / r->num << " usec average, throughput "
                  << double(r->num) * valueSize / r->usec << " MB/s" << std::endl;
    }
    gettimeofday(&tv, 0);
    long end = tv.tv_sec * 1000000 + tv.tv_usec;
    std::cout << total << " total records get in " << end - start << " usec,"
              << double(end - start) / total << " usec average, throughput "
              << double(total) / (end - start) * valueSize << " MB/s, qps is "
              << double(total) / (end - start) * 1000000 << std::endl;
}
