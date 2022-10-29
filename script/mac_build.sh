#!/bin/sh

dir_path=$(pwd)
BINARY_FILE=$(basename "$dir_path")
BUILD_DIR=${dir_path}/build

echo "$BINARY_FILE"
echo "$BUILD_DIR"

cd src && go build -o "${BUILD_DIR}"/"${BINARY_FILE}"
cd - || exit 1

#创建一个build文件夹，用来保存编译的结果
#将本地的etc下的文件拷贝只build/etc下
mkdir -p "${BUILD_DIR}"/etc
cp -r ./etc/* "${BUILD_DIR}"/etc

#拷贝启动脚本
cp ./script/start.sh "${BUILD_DIR}"
