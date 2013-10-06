#!/bin/sh

set -eux

cd regexptransform/
go test
cd ..

cd tftp
go test
cd ..

cd config
go test
cd ..

cd hooks
go test
cd ..

go test
