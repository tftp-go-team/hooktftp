#!/bin/sh

set -eu
set -x

sudo apt-get update
sudo apt-get install wget puavo-devscripts gcc --yes --force-yes

# Too lazy to backport Go 1.2 for Ubuntu Precise. We'll just install the binary
# manually.

cd .. # get out of the source checkout

arch=
processor="$(gcc -dumpmachine)"

[ "$processor" = "i686-linux-gnu" ] && arch="386"
[ "$processor" = "x86_64-linux-gnu" ] && arch="amd64"
if [ "$arch" = "" ]; then
    echo "Unknown processor $processor"
    exit 1
fi

go_bin_name="go1.2.linux-${arch}.tar.gz"
wget -c "https://go.googlecode.com/files/${go_bin_name}"
tar xzvf go*.tar.gz

export GOROOT="$(pwd)/go"
export PATH="$PATH:$GOROOT/bin"



# Setup Go workspace for compiling
mkdir workspace
cd workspace

export GOPATH="$(pwd)"

# Move sources to the correct workspace location
mkdir -p src/github.com/epeli
mv ../hooktftp src/github.com/epeli
cd src/github.com/epeli/hooktftp

# Do Opinsys style debian package build
puavo-build-debian-dir
sudo puavo-install-deps debian/control
puavo-dch $(cat VERSION)
puavo-debuild

# Upload packages to apt repositories
aptirepo-upload -r $APTIREPO_REMOTE -b "git-$(echo "$GIT_BRANCH" | cut -d / -f 2)" ../hooktftp*.changes
aptirepo-upload -r $APTIREPO_REMOTE -b hooktftp ../hooktftp*.changes
