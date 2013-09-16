#!/bin/bash

OUTDIR="out$1"

set -eu

fetch() {
    echo "Fetching $1"
    atftp --get --remote-file fixtures/$1 --local-file $OUTDIR/$1 localhost 1234 || true
}

if [ ! -d fixtures ]; then
    echo "Fixtures missing! Run tools/create-fixtures.sh"
    exit 1
fi

rm -rf $OUTDIR
mkdir $OUTDIR

echo "Fetching files"
fetch small &
fetch medium &
fetch mod512 &
fetch mod512double &
fetch big &
atftp --option "blksize 100" --get --remote-file fixtures/medium2 --local-file $OUTDIR/medium2 localhost 1234 &
atftp --option "blksize 1536" --get --remote-file fixtures/big2 --local-file $OUTDIR/big2 localhost 1234 &

wait

cd $OUTDIR
sha1sum --check ../fixtures/SHA1SUMS
cd ..

set +e
ERROR_MESSAGE=$(atftp --get --remote-file nonexistent --local-file /dev/null localhost 1234 2>&1)
set -e
if [[ ! $ERROR_MESSAGE =~ "no such file or directory" ]]; then
    echo "Cannot find 'no such file or directory' from: $ERROR_MESSAGE"
    exit 1
fi

set +e
ERROR_MESSAGE=$(atftp --get --remote-file ../foo.txt --local-file /dev/null localhost 1234 2>&1)
set -e
if [[ ! $ERROR_MESSAGE =~ "Path access violation" ]]; then
    echo "Cannot find 'Path access violation' from: $ERROR_MESSAGE"
    exit 1
fi

atftp --get --remote-file custom.txt --local-file $OUTDIR/custom.txt localhost 1234
CONTENT=$(cat $OUTDIR/custom.txt)
if [ "$CONTENT" != "customdata" ]; then
    echo "Did not receive custom data for custom.txt"
    exit 1
fi


echo "ALL OK"
