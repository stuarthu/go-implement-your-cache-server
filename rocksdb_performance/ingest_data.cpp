#include "common.h"

#include <sys/time.h>

int main(int argc, char** argv)
{
    // ingest total * valueSize = 10GB data
    variables_map vm = parse(argc, argv, 0);
    int total = vm["total"].as<int>();
    int valueSize = vm["size"].as<int>();

    DB* db = opendb();

    struct timeval tv;
    gettimeofday(&tv, 0);
    long start = tv.tv_sec * 1000000 + tv.tv_usec;

    std::string valuePrefix = std::string(valueSize, 'a');
    for (int i = 0; i < total; i++) {
        std::string key = "ingest" + std::to_string(i);
        std::string value = valuePrefix + std::to_string(i);
        db->Put(WriteOptions(), key, value);
    }

    gettimeofday(&tv, 0);
    long end = tv.tv_sec * 1000000 + tv.tv_usec;
    std::cout << total << " total records set in " << end - start << " usec,"
              << double(end - start) / total << " usec average, throughput "
              << double(total) / (end - start) * valueSize << " MB/s, rps is "
              << double(total) / (end - start) * 1000000 << std::endl;
}
