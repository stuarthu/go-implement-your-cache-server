#include "common.h"

#include <sys/time.h>

int main(int argc, char** argv)
{
    int key;
    options_description keyOption("key option");
    keyOption.add_options()("key,k", value<int>(&key), "key");

    variables_map vm = parse(argc, argv, &keyOption);
    int total = vm["total"].as<int>();
    int valueSize = vm["size"].as<int>();
    int s = 0;
    int e = total;
    if (key) {
        std::cout << "key is " << key << std::endl;
        s = key;
        e = key + 1;
    }

    DB* db = opendb(true);

    std::string valuePrefix = std::string(valueSize, 'a');
    std::string value;
    int fail=0;
    struct timeval tv;
    gettimeofday(&tv, 0);
    int start = tv.tv_sec * 1000000 + tv.tv_usec;
    for (int i = s; i < e; i++) {
        std::string tmp = std::to_string(i);
        std::string key = keyPrefix + tmp;
        auto s = db->Get(ReadOptions(), key, &value);
        if (!s.ok()) {
            std::cerr << "db->Get(): key=" << key
                << "," << s.ToString() << std::endl;
            fail++;
        }
        else {
            assert(value == valuePrefix + tmp);
        }
    }
    gettimeofday(&tv, 0);
    int end = tv.tv_sec * 1000000 + tv.tv_usec;
    std::cout << total << " records get in " << end - start << " usec, "
              << double(end - start) / total << " usec average, fail "
              << fail << ", fail rate " << 100*fail/total << "%" << std::endl;

    delete db;
}
