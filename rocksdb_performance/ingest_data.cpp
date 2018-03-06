#include "common.h"

int main(int argc, char** argv)
{
    DB* db = opendb();

    // ingest total * valueSize = 10GB data
    int total = 100000;
    int valueSize = 100000;

    std::string valuePrefix = std::string(valueSize, 'a');
    for (int i = 0; i < total; i++) {
        std::string key = "ingest" + std::to_string(i);
        std::string value = valuePrefix + std::to_string(i);
        db->Put(WriteOptions(), key, value);
    }

    delete db;
}
