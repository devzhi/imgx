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
      专为解决Docker Hub拉取镜像困难而设计的效率工具
    </b>
  </p>
  <p>
  <a href="https://github.com/devzhi/MessageRouter/blob/main/LICENSE"><img src="https://img.shields.io/badge/license-APL2.0-blue.svg"></img></a>
  <a href="#"><img src="https://img.shields.io/badge/Contributions-welcome-green?logo=github"></img></a>
    <a href="#"><img src="https://img.shields.io/badge/Version-0.0.1-green.svg"></img></a>
  </p>
</div>

## 项目愿景

在网络状况良好的设备中拉取镜像并自动推送至目标服务器，以解决Docker Hub拉取镜像困难的问题。

## 功能特性

- [ ] 从Docker Hub拉取镜像
- [x] 在不依赖Docker的情况下从Docker Hub拉取镜像
- [ ] 自动推送镜像至目标服务器

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
imgx [options]
```

### 选项说明

`-arch string`：用于指定镜像（image）的架构，默认值为 "amd64"。例如，若需指定为其他架构（如 arm64），则添加 `-arch arm64` 选项。

`-name string`：用来设定镜像（image）的名称，根据实际需求填写相应的名称内容即可。

`-os string`：指定镜像（image）的操作系统类型，默认值是 "linux"。若针对其他操作系统的镜像操作，可添加 `-os` 选项。

`-tag string`：确定镜像（image）的标签，默认标签为 "latest"。如需自定义标签（比如 v1.0），可添加 `-tag v1.0` 选项。

### 示例

```shell
imgx -name nginx
imgx -name nginx -tag latest
imgx -name nginx -tag latest -arch arm64
imgx -name nginx -tag latest -arch arm64 -os linux
```

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=devzhi/imgx&type=Date)](https://star-history.com/#devzhi/imgx&Date)

## 免责声明

本开源项目（imgx）是依据 [Apache License 2.0](https://github.com/devzhi/MessageRouter/blob/main/LICENSE) 许可证授权发布，旨在为开发者社区提供有益的工具、代码或资源，以促进技术创新与共享。

使用者应知悉并同意：

1. 本项目仅用于合法合规的目的，严禁将本项目的任何部分用于任何非法活动，包括但不限于未经授权的访问、数据窃取、恶意软件传播、侵犯他人知识产权或违反法律法规的行为。若使用者违反此规定，造成的一切后果与法律责任将由使用者自行承担，与本项目的开发者及贡献者无关。
2. 本项目按“原样”提供，不提供任何形式的明示或暗示的保证，包括但不限于对项目的适用性、准确性、可靠性、完整性以及不侵权的保证。使用者在使用本项目时应自行承担风险，开发者及贡献者不对因使用本项目而产生的任何直接或间接损失负责。

请在使用本开源项目前仔细阅读并理解本免责声明，若您选择使用本项目，即表示您已接受上述条款。