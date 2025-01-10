#!/bin/bash

# 使用协议
echo "使用协议"
echo "本开源项目（imgx）是依据 APL 2.0 许可证授权发布，旨在为开发者社区提供有益的工具、代码或资源，以促进技术创新与共享。"
echo "使用者应知悉并同意："
echo "1. 本项目仅用于合法合规的目的，严禁将本项目的任何部分用于任何非法活动，包括但不限于未经授权的访问、数据窃取、恶意软件传播、侵犯他人知识产权或违反法律法规的行为。若使用者违反此规定，造成的一切后果与法律责任将由使用者自行承担，与本项目的开发者及贡献者无关。"
echo "2. 本项目按“原样”提供，不提供任何形式的明示或暗示的保证，包括但不限于对项目的适用性、准确性、可靠性、完整性以及不侵权的保证。使用者在使用本项目时应自行承担风险，开发者及贡献者不对因使用本项目而产生的任何直接或间接损失负责。"
echo "请在使用本开源项目前仔细阅读并理解本协议，若您选择使用本项目，即表示您已接受上述条款。"
echo
read -p "是否接受使用协议？(Y/y|N/n): " accept < /dev/tty

if [[ "$accept" != [Yy] ]]; then
    echo "您未接受使用协议，安装终止。"
    exit 1
fi

# 检查/usr/local/bin目录是否可写
if [ ! -w /usr/local/bin ]; then
    echo "/usr/local/bin 目录不可写，请使用 sudo 或 root 用户运行此脚本"
    exit 1
fi

# 检查/opt目录是否可写
if [ ! -w /opt ]; then
    echo "/opt 目录不可写，请使用 sudo 或 root 用户运行此脚本"
    exit 1
fi

# 检查/opt/imgx目录是否存在，不存在则创建
if [ ! -d /opt/imgx ]; then
    mkdir -p /opt/imgx
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
ln -sf /opt/imgx/imgx /usr/local/bin/imgx

if [ $? -eq 0 ]; then
    echo "imgx ${latest_version} 安装成功"
else
    echo "imgx ${latest_version} 安装失败"
    exit 1
fi