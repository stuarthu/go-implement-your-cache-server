#include "common.h"

#include <sys/time.h>

int main(int argc, char** argv)
{
    int batchSize;
    options_description batchOption("batch_write option");
    batchOption.add_options()("batch_size,b", value<int>(&batchSize)->default_value(1), "batch size");

    variables_map vm = parse(argc, argv, &batchOption);
    int total = vm["total"].as<int>();
    int valueSize = vm["size"].as<int>();
    std::cout << "batch size is " << batchSize << std::endl;

    DB* db = opendb();

    struct timeval tv;
    gettimeofday(&tv, 0);
    long start = tv.tv_sec * 1000000 + tv.tv_usec;
    int count = 0;
    std::string valuePrefix = std::string(valueSize, 'a');
    while (count != total) {
        WriteBatch batch;
        for (int j = 0; j < batchSize; j++, count++) {
            if (count == total)
                break;
            std::string tmp = std::to_string(count);
            std::string key = keyPrefix + tmp;
            std::string value = valuePrefix + tmp;
            auto s = batch.Put(key, value);
            if (!s.ok())
                std::cerr << "batch.Put():" << s.ToString() << std::endl;
            assert(s.ok());
        }
        auto s = db->Write(WriteOptions(), &batch);
        if (!s.ok())
            std::cerr << "db->Write():" << s.ToString() << std::endl;
        assert(s.ok());
    }
    gettimeofday(&tv, 0);
    long end = tv.tv_sec * 1000000 + tv.tv_usec;
    std::cout << total << " records batch put in " << end - start << " usec, "
              << double(end - start) / total << " usec average, throughput is "
              << (double)total * valueSize / (end - start) << " MB/s, qps is "
              << (double)1000000 * total / (end - start) << std::endl;

    delete db;
}
