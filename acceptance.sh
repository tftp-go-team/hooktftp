#!/bin/sh

set -eu


fetch() {
    echo "Fetching $1"
    atftp --get --remote-file fixtures/$1 --local-file out/$1 localhost 1234 || true
}


if [ ! -d fixtures ]; then
    echo "Creating fixtures"
    mkdir fixtures
    cd fixtures
    echo "smalfile" > small
    dd if=/dev/urandom of=medium bs=1048577 count=5
    dd if=/dev/urandom of=medium2 bs=1048577 count=5
    dd if=/dev/urandom of=big bs=1048577 count=10
    dd if=/dev/urandom of=big2 bs=1048577 count=10
    dd if=/dev/urandom of=mod512 bs=512 count=1
    dd if=/dev/urandom of=mod512double bs=512 count=2
    echo "Writing checks sums"
    sha1sum * > SHA1SUMS
    cd ../
fi

rm -rf out
mkdir out

echo "Fetching files"
fetch small
fetch medium
fetch mod512
fetch mod512double
fetch big
atftp --option "blksize 100" --get --remote-file fixtures/medium2 --local-file out/medium2 localhost 1234
atftp --option "blksize 1536" --get --remote-file fixtures/big2 --local-file out/big2 localhost 1234

cd out
sha1sum --check ../fixtures/SHA1SUMS
