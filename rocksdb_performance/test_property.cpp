#include "common.h"

#include <sys/time.h>

int main(int argc, char** argv)
{
    DB* db = opendb();
    std::string pValue;
    TablePropertiesCollection c;
    db->GetPropertiesOfAllTables(&c);
    int count = 0;
    for (TablePropertiesCollection::iterator i = c.begin(); i != c.end(); i++) {
        std::cout << i->first << i->second->raw_key_size << std::endl;
        count += i->second->raw_key_size;
    }
    db->GetProperty("rocksdb.aggregated-table-properties", &pValue);
    std::cout << pValue << std::endl;
    std::cout << count << std::endl;

    delete db;
}
