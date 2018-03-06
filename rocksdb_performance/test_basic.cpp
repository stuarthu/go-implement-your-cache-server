#include "common.h"

#include <sys/time.h>

int main(int argc, char** argv)
{
    variables_map vm = parse(argc, argv, 0);
    int total = vm["total"].as<int>();
    int valueSize = vm["size"].as<int>();

    DB* db = opendb();

    struct timeval tv;
    gettimeofday(&tv, 0);
    long start = tv.tv_sec * 1000000 + tv.tv_usec;
    std::string valuePrefix = std::string(valueSize, 'a');
    for (int i = 0; i < total; i++) {
        std::string key = keyPrefix + std::to_string(i);
        std::string value = valuePrefix + std::to_string(i);
        db->Put(WriteOptions(), key, value);
    }
    gettimeofday(&tv, 0);
    long end = tv.tv_sec * 1000000 + tv.tv_usec;
    std::cout << total << " records put in " << end - start << " usec, "
              << double(end - start) / total << " usec average" << std::endl;

    gettimeofday(&tv, 0);
    std::string value;
    start = tv.tv_sec * 1000000 + tv.tv_usec;
    for (int i = 0; i < total; i++) {
        std::string tmp = std::to_string(std::rand() % total);
        std::string key = keyPrefix + tmp;
        db->Get(ReadOptions(), key, &value);
        assert(value == valuePrefix + tmp);
    }
    gettimeofday(&tv, 0);
    end = tv.tv_sec * 1000000 + tv.tv_usec;
    std::cout << total << " records get in " << end - start << " usec, "
              << double(end - start) / total << " usec average" << std::endl;

    delete db;
}
