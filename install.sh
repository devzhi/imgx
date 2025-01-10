#!/bin/bash

# 获取最新版本
latest_version=$(curl -s https://api.github.com/repos/devzhi/imgx/releases/latest | grep tag_name | cut -d '"' -f 4)

# 获取系统类型
os=$(uname -s)

# 获取系统架构
arch=$(uname -m)

# 下载地址
download_url="https://github.com/devzhi/imgx/releases/download/${latest_version}/imgx_${latest_version}_${os}_${arch}.tar.gz"

# 下载并安装
echo "正在下载 imgx ${latest_version} ..."
curl -sL ${download_url} | tar xz -C /usr/local/bin

if [ $? -eq 0 ]; then
    echo "imgx ${latest_version} 安装成功"
else
    echo "imgx ${latest_version} 安装失败"
    exit 1
fi