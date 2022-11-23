#!/usr/bin/env bash

VERSION=$(grep -Po '_Version = "(\d+\.){2}\d+"' gosvr.go | cut -d\" -f2)
make cleanall
mkdir release
declare -a arr=("darwin_386" "darwin_amd64" "darwin_arm64" "linux_386" "linux_amd64" "linux_arm" "linux_arm64" "windows_386" "windows_amd64")
for i in "${arr[@]}"
do
    packr2 clean
    OS=$(cut -d'_' -f1 <<<"$i")
    ARCH=$(cut -d'_' -f2 <<<"$i")
    echo "Building $OS $ARCH ..."
    env GOOS=$OS GOARCH=$ARCH packr2
    cd release
    env GOOS=$OS GOARCH=$ARCH go build github.com/JeziL/gosvr
    tar -czf b_gosvr_v${VERSION}_${OS}_$ARCH.tar.gz ./gosvr*  --remove-files
    cd ..
done
cd release
rename 's/b_//' *.tar.gz
cd ..
packr2 clean
