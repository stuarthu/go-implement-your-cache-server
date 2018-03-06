#include <iostream>

#include <boost/program_options.hpp>

#include "rocksdb/db.h"
#include "rocksdb/utilities/db_ttl.h"

using namespace rocksdb;
using namespace boost::program_options;

const std::string keyPrefix = std::string(99, 'a');

variables_map parse(int argc, char** argv, options_description* additional)
{
    int total, valueSize;
    options_description desc("Allowed options");
    desc.add_options()("help,h", "produce help message")("total,t", value<int>(&total)->default_value(10000), "total record number")("size,s", value<int>(&valueSize)->default_value(1000), "value size");

    if (additional) {
        desc.add(*additional);
    }

    variables_map vm;
    store(parse_command_line(argc, argv, desc), vm);
    notify(vm);

    if (vm.count("help")) {
        std::cout << desc << std::endl;
        exit(1);
    }

    std::cout << "total record number is " << total << std::endl;
    std::cout << "value size is " << valueSize << std::endl;

    return vm;
}

DB* opendb(bool readonly = false, const std::string& dir = "/mnt/rocksdb", int ttl = 0)
{
    DB* db;
    Options options;
    // Optimize RocksDB. This is the easiest way to get RocksDB to perform well
    options.IncreaseParallelism();
    options.OptimizeLevelStyleCompaction();
    // create the DB if it's not already present
    options.create_if_missing = true;

    // open DB
    Status s;
    if (readonly)
        s = DB::OpenForReadOnly(options, dir, &db);
    else if (ttl) {
        rocksdb::DBWithTTL* db_ttl;
        s = rocksdb::DBWithTTL::Open(options, dir, &db_ttl, ttl);
        db = db_ttl;
    }
    else
        s = DB::Open(options, dir, &db);

    if (!s.ok()) {
        std::cout << "open " << dir << ":" << s.ToString() << std::endl;
    }
    assert(db);
    return db;
}
