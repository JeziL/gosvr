#!/usr/bin/env bash

VERSION=$(grep -Po '_Version = "(\d+\.){2}\d+"' gosvr.go | cut -d\" -f2)
make cleanall
mkdir release
declare -a arr=("darwin_386" "darwin_amd64" "dragonfly_amd64" "freebsd_386" "freebsd_amd64" "freebsd_arm" "linux_386" "linux_amd64" "linux_arm" "linux_arm64" "linux_ppc64" "linux_ppc64le" "linux_mips" "linux_mipsle" "linux_mips64" "linux_mips64le" "netbsd_386" "netbsd_amd64" "netbsd_arm" "openbsd_386" "openbsd_amd64" "openbsd_arm" "plan9_386" "plan9_amd64" "solaris_amd64" "windows_386" "windows_amd64")
for i in "${arr[@]}"
do
    packr clean
    OS=$(cut -d'_' -f1 <<<"$i")
    ARCH=$(cut -d'_' -f2 <<<"$i")
    echo "Building $OS $ARCH ..."
    env GOOS=$OS GOARCH=$ARCH packr
    cd release
    env GOOS=$OS GOARCH=$ARCH go build github.com/JeziL/gosvr
    tar -czf b_gosvr_v${VERSION}_${OS}_$ARCH.tar.gz ./gosvr*  --remove-files
    cd ..
done
cd release
rename 's/b_//' *.tar.gz
cd ..
packr clean