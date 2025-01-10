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