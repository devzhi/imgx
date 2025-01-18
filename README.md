<div align="center">
  <p>
    <h3>
      <b>
        imgx
      </b>
    </h3>
  </p>
  <p>
    <b>
      Docker镜像传输工具
    </b>
  </p>
  <p>
  <a href="https://github.com/devzhi/imgx/blob/main/LICENSE"><img src="https://img.shields.io/badge/license-APL2.0-blue.svg"></img></a>
  <a href="https://github.com/devzhi/imgx/pulls"><img src="https://img.shields.io/badge/Contributions-welcome-green?logo=github"></img></a>
  <a href="https://github.com/devzhi/imgx/releases/latest"><img src="https://img.shields.io/badge/Version-1.0.2-green.svg"></img></a>
  </p>
</div>

## 功能特性

- [x] 从Docker Hub拉取镜像（不依赖Docker）
- [x] 将本地镜像文件推送至远程主机
- [x] 从Docker Hub拉取镜像并推送至远程主机

## 使用说明

### 安装

#### Linux / macOS

```shell
curl https://raw.githubusercontent.com/devzhi/imgx/main/install.sh | sudo bash
```

#### Windows

1. 进入[Release](https://github.com/devzhi/imgx/releases/latest)页面下载最新版本的imgx压缩包
2. 解压压缩包
3. 将解压后的imgx可执行文件添加至系统环境变量中

### 调用方式

```shell
imgx [flags]
imgx [command]
```

### 可用命令

`completion`：生成指定shell的自动补全脚本

`help`：关于任何命令的帮助

`load`：将镜像加载到远程主机

`pull`：从Docker hub本地拉取镜像

`version`：显示imgx版本信息

`x`：拉取并加载镜像到远程主机

### Flags 标志

`-h, --help`：帮助信息
`-v, --version`：显示版本信息

### 示例

#### 从Docker Hub拉取镜像

```shell
imgx pull [image] [flags]

Flags:
  -a, --arch string   拉取镜像的架构 (默认 "amd64") [可选]
  -h, --help          帮助信息 [可选]
  -o, --os string     拉取镜像的操作系统 (默认 "linux") [可选]
  -t, --tag string    拉取镜像的标签 (默认 "latest") [可选]

imgx pull nginx -a amd64 -o linux -t latest
```

#### 将镜像加载到远程主机

```shell
imgx load [input] [flags]

Flags:
  -h, --help              帮助信息 [可选]
  -H, --host string       远程主机地址
  -p, --password bool     使用密码登录远程主机（默认 false）[可选]
  -P, --port int          远程主机的端口 (默认 22) [可选]
      --protocol string   远程主机的SSH协议 (默认 "tcp") [可选]
  -u, --username string   远程主机的用户名

imgx load nginx_latest_amd64_linux.tar.gz -H 192.168.1.100 -P 22 -u user -p password --protocol tcp -r
```

#### 拉取并加载镜像到远程主机

```shell
imgx x [image] [flags]

Flags:
  -a, --arch string       拉取镜像的架构 (默认 "amd64") [可选]
  -h, --help              帮助信息 [可选]
  -H, --host string       远程主机地址
  -o, --os string         拉取镜像的操作系统 (默认 "linux") [可选]
  -p, --password bool     使用密码登录远程主机（默认 false）[可选]
  -P, --port int          远程主机的端口 (默认 22) [可选]
      --protocol string   远程主机的SSH协议 (默认 "tcp") [可选]
  -s, --save              成功加载后保留镜像文件（默认 "false"） [可选]
  -t, --tag string        拉取镜像的标签 (默认 "latest") [可选]
  -u, --username string   远程主机的用户名

imgx x nginx -a amd64 -o linux -t latest -H 192.168.1.100 -P 22 -u user -p password --protocol tcp -r
```

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=devzhi/imgx&type=Date)](https://star-history.com/#devzhi/imgx&Date)

## 免责声明

本开源项目（imgx）是依据 [Apache License 2.0](https://github.com/devzhi/imgx/blob/main/LICENSE)
许可证授权发布，旨在为开发者社区提供有益的工具、代码或资源，以促进技术创新与共享。

使用者应知悉并同意：

1.
本项目仅用于合法合规的目的，严禁将本项目的任何部分用于任何非法活动，包括但不限于未经授权的访问、数据窃取、恶意软件传播、侵犯他人知识产权或违反法律法规的行为。若使用者违反此规定，造成的一切后果与法律责任将由使用者自行承担，与本项目的开发者及贡献者无关。
2. 本项目按“原样”提供，不提供任何形式的明示或暗示的保证，包括但不限于对项目的适用性、准确性、可靠性、完整性以及不侵权的保证。使用者在使用本项目时应自行承担风险，开发者及贡献者不对因使用本项目而产生的任何直接或间接损失负责。

请在使用本开源项目前仔细阅读并理解本免责声明，若您选择使用本项目，即表示您已接受上述条款。