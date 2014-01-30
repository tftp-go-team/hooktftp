#!/bin/sh

set -eu
set -x

# https://go.googlecode.com/files/go1.2.linux-amd64.tar.gz
# https://go.googlecode.com/files/go1.2.linux-386.tar.gz

sudo apt-get update
sudo apt-get install wget puavo-devscripts --yes --force-yes

arch=
processor="$(uname --processor)"

[ "$processor" = "i686" ] && arch="386"
[ "$processor" = "x86_64" ] && arch="amd64"
if [ "$arch" = "" ]; then
    echo "Unknown processor $processor"
    exit 1
fi

cd ..

go_bin_name="go1.2.linux-${arch}.tar.gz"

wget -c "https://go.googlecode.com/files/${go_bin_name}"
tar xzvf go*.tar.gz

export GOROOT="$(pwd)/go"
export PATH="$PATH:$GOROOT/bin"

mkdir workspace
cd workspace

export GOPATH="$(pwd)"

mkdir -p src/github.com/epeli
mv ../hooktftp src/github.com/epeli
cd src/github.com/epeli/hooktftp

puavo-build-debian-dir
sudo puavo-install-deps debian/control
puavo-dch $(cat VERSION)
puavo-debuild

aptirepo-upload -r $APTIREPO_REMOTE -b "git-$(echo "$GIT_BRANCH" | cut -d / -f 2)" ../hooktftp*.changes
