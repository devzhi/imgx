#!/bin/bash

# 检查/usr/local/bin目录是否可写
if [ ! -w /usr/local/bin ]; then
    echo "/usr/local/bin 目录不可写，请使用 sudo 或 root 用户运行此脚本"
    exit 1
fi

# 检查是否安装了 curl
if ! command -v curl &> /dev/null; then
    echo "请先安装 curl"
    exit 1
fi

# 获取最新版本
latest_version=$(curl -s https://api.github.com/repos/devzhi/imgx/releases/latest | grep tag_name | cut -d '"' -f 4)

# 获取系统类型
os=$(uname -s)

# 获取系统架构
arch=$(uname -m)

# 下载地址
download_url="https://github.com/devzhi/imgx/releases/download/${latest_version}/imgx_${latest_version#v}_${os}_${arch}.tar.gz"

# 下载并安装到/opt/imgx
echo "正在下载 imgx ${latest_version} ..."
curl -sL ${download_url} | tar xzf - -C /opt/imgx/

if [ $? -ne 0 ]; then
    echo "imgx ${latest_version} 安装失败"
    exit 1
fi

# 创建软连接到/usr/local/bin
sudo ln -sf /opt/imgx/imgx /usr/local/bin/imgx

if [ $? -eq 0 ]; then
    echo "imgx ${latest_version} 安装成功"
else
    echo "imgx ${latest_version} 安装失败"
    exit 1
fi