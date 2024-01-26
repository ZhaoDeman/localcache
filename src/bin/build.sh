#!/bin/bash
set -x
current_path=$PWD
cd current_path

cd .././main && go build -o localcache
./localcache -config_file ../config.json