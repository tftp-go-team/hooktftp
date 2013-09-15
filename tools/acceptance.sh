#!/bin/sh

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
fetch small
fetch medium
fetch mod512
fetch mod512double
fetch big
atftp --option "blksize 100" --get --remote-file fixtures/medium2 --local-file $OUTDIR/medium2 localhost 1234
atftp --option "blksize 1536" --get --remote-file fixtures/big2 --local-file $OUTDIR/big2 localhost 1234

cd $OUTDIR
sha1sum --check ../fixtures/SHA1SUMS
