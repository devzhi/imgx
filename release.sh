#!/bin/bash

# 检查是否提供了版本号
if [ -z "$1" ]; then
    echo "Usage: $0 <version>"
    exit 1
fi

# 提取版本号
version=$1
tag="v$version"

# 更新 main.go 中的版本号
sed -i "s/version := \"[^\"]*\"/version := \"$version\"/" ./cmd/version.go

# 更新 README.md 中的版本号
sed -i "s/Version-\([0-9]\{1,\}\.\)\{2\}[0-9]\{1,\}/Version-$version/g" ./README.md

# 添加修改后的文件到暂存区
git add ./main.go
git add ./README.md

# 提交修改
git commit -m "Update version to $version"

# 创建新的 tag
git tag $tag

# 推送提交和 tag
git push origin main
git push origin $tag

echo "Version $version released successfully."