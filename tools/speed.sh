#!/bin/sh

if [ ! -d fixtures ]; then
    echo "Fixtures missing! Run tools/create-fixtures.sh"
    exit 1
fi

for i in $(seq 10)
do
    atftp --option "blksize 1536" --get --remote-file fixtures/big2 --local-file /dev/null localhost 1234
done
