#!/bin/bash

# 设置目标平台
platforms=("linux/amd64" "linux/arm64" "windows/amd64" "darwin/amd64" "darwin/arm64")

# 设置版本
version="v0.2.1"

# 设置构建时间
build_time=$(date +'%Y%m%d')
#build_time="20241114"

# 设置路径
path="G:/gitlab/Subfinder/cmd/main.go"

# 设置输出目录
output_dir="build"

# 创建输出目录
mkdir -p $output_dir

# 遍历每个平台进行编译
for platform in "${platforms[@]}"
do
    # 分割平台字符串
    IFS="/" read -r -a platform_split <<< "$platform"
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}

    # 设置输出文件名
    output_name="$output_dir/Subfinder_${GOOS}_${GOARCH}_${version}_${build_time}"
    if [ "$GOOS" = "windows" ]; then
        output_name+=".exe"
    fi

    # 编译
    echo "Building for $platform..."
    env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name $path

    if [ $? -ne 0 ]; then
        echo "An error has occurred! Aborting the script execution..."
        read -p "按任意键退出..."

    fi
done

echo "Build completed!"
read -p "按任意键退出..."
