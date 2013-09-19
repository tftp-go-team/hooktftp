#!/bin/bash

OUTDIR="out$1"

set -eu

fetch() {
    echo "Fetching $1"
    atftp --get --remote-file fixtures/$1 --local-file $OUTDIR/$1 localhost 1234 || true
}

contains() {
    echo
    echo
    if [[ "$1" == *"$2"* ]]; then
        echo "OK '$2'"
        echo
    else
        echo "FAIL: Did not find '$2' from '$1'"
        exit 1
    fi

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
ERROR_MESSAGE=$(atftp --get --remote-file small --local-file /dev/null localhost 1234 2>&1)
set -e
contains "$ERROR_MESSAGE" "no such file or directory"

set +e
ERROR_MESSAGE=$(atftp --get --remote-file ../foo.txt --local-file /dev/null localhost 1234 2>&1)
set -e
contains "$ERROR_MESSAGE" "Path access violation"

atftp --get --remote-file custom.txt --local-file $OUTDIR/custom.txt localhost 1234
CONTENT=$(cat $OUTDIR/custom.txt)
if [ "$CONTENT" != "customdata" ]; then
    echo "FAIL Did not receive custom data for custom.txt"
    exit 1
else
    echo "OK custom data"
fi


echo "ALL OK"
