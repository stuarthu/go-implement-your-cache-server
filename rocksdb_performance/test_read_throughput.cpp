#include "common.h"

#include <sys/time.h>

struct threadOption {
    DB* db;
    int num;
    int valueSize;
    int total;
};

struct result {
    int num;
    long usec;
    int fail;
};

void* func(void* arg)
{
    threadOption* o = (threadOption*)arg;
    DB* db = o->db;
    int num = o->num;
    int fail = 0;
    std::string value;
    std::string valuePrefix = std::string(o->valueSize, 'a');
    struct timeval tv;
    gettimeofday(&tv, 0);
    long start = tv.tv_sec * 1000000 + tv.tv_usec;
    for (int i = 0; i < num; i++) {
        std::string tmp = std::to_string(rand() % o->total);
        std::string key = keyPrefix + tmp;
        auto s = db->Get(ReadOptions(), key, &value);
        if (!s.ok()) {
            std::cerr << "db->Get(): key=" << key
                      << "," << s.ToString() << std::endl;
            fail++;
        } else {
            assert(value == valuePrefix + tmp);
        }
    }
    gettimeofday(&tv, 0);
    long end = tv.tv_sec * 1000000 + tv.tv_usec;
    result* r = new result;
    r->num = num;
    r->usec = end - start;
    r->fail = fail;
    return r;
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

    DB* db = opendb(true);

    pthread_t get_thread[threads];
    struct timeval tv;
    gettimeofday(&tv, 0);
    long start = tv.tv_sec * 1000000 + tv.tv_usec;
    for (int i = 0; i < threads; i++) {
        threadOption* o = new threadOption;
        o->num = total / threads;
        o->db = db;
        o->valueSize = valueSize;
        o->total = total;
        pthread_create(&get_thread[i], 0, func, o);
    }

    int fail = 0;
    for (int i = 0; i < threads; i++) {
        void* res;
        pthread_join(get_thread[i], &res);
        result* r = (result*)res;
        fail += r->fail;
        std::cout << r->num << " records get in " << r->usec << " usec, "
                  << double(r->usec) / r->num << " usec average, throughput "
                  << double(r->num) * valueSize / r->usec << " MB/s, fail "
                  << r->fail << ", fail rate " << 100*r->fail/r->num << "%" << std::endl;
    }
    gettimeofday(&tv, 0);
    long end = tv.tv_sec * 1000000 + tv.tv_usec;
    std::cout << total << " total records get in " << end - start << " usec,"
              << double(end - start) / total << " usec average, throughput "
              << double(total) / (end - start) * valueSize << " MB/s, qps is "
              << double(total) / (end - start) * 1000000 << ", fail " << fail
              << ", fail rate " << 100*fail/total << "%" << std::endl;

    delete db;
}
