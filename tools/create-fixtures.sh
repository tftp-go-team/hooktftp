#!/bin/sh

echo "Creating fixtures"
rm -rf fixtures
mkdir fixtures
cd fixtures
echo "smalfile" > small
dd if=/dev/urandom of=medium bs=1048577 count=5
dd if=/dev/urandom of=medium2 bs=1048577 count=5
dd if=/dev/urandom of=webfile bs=1048577 count=5
dd if=/dev/urandom of=big bs=1048577 count=10
dd if=/dev/urandom of=big2 bs=1048577 count=50
dd if=/dev/urandom of=mod512 bs=512 count=1
dd if=/dev/urandom of=mod512double bs=512 count=2
echo "Writing checks sums"
sha1sum * > SHA1SUMS
cd ../
