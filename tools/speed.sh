#!/bin/sh

if [ ! -d fixtures ]; then
    echo "Fixtures missing! Run tools/create-fixtures.sh"
    exit 1
fi

name="testfile-$(shuf -i 2000-65000 -n 1)"

for i in $(seq 10)
do
    echo "Starting $name-$i"
    atftp --get --remote-file fixtures/big --local-file "$name-$i" localhost 1234
    echo "Ending $name-$i"
done
