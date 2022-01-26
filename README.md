# Articli

[![Release](https://github.com/k8scat/Articli/actions/workflows/release.yaml/badge.svg)](https://github.com/k8scat/Articli/actions/workflows/release.yaml)
[![GitHub Repo stars](https://img.shields.io/github/stars/k8scat/articli?style=social)](https://github.com/k8scat/Articli/stargazers)
[![GitHub watchers](https://img.shields.io/github/watchers/k8scat/articli?style=social)](https://github.com/k8scat/Articli/watchers)
[![star](https://gitee.com/k8scat/articli/badge/star.svg?theme=dark)](https://gitee.com/k8scat/articli/stargazers)
[![codecov](https://codecov.io/gh/k8scat/Articli/branch/main/graph/badge.svg?token=045FCRVF27)](https://codecov.io/gh/k8scat/Articli)

**Articli** is an Article CLI tool for managing content in multi platforms.

**Articli** 是一个可以管理多个平台内容的命令行工具，
通过解析 `Markdown` 文件内容以及调用平台接口，实现内容管理。

最终目标是基于 **本地文件** + **Git 代码仓** 管理所有的文章，
并且可以通过命令行操作以及 CI/CD，实现文章在各个平台的发布、更新等功能。
这样做的好处有：

- 数据安全，既发布到了第三方平台，又可以通过 **Git 代码仓**管理，避免因平台问题导致数据丢失
- 可以实现自动化，比如文章自动在多个平台发布、更新
- 面向程序员的 CLI 工具，可以实现更多个性化的操作

为本项目点赞将鼓励作者继续完善下去，欢迎提出建议、Bug、PR。

## 支持的平台

- [掘金](https://juejin.cn)
- [开源中国](https://oschina.net)
- [GitHub](https://github.com)

## 安装

### NPM

```shell
npm install -g @k8scat/articli
```

### Homebrew

```shell
# 添加 tap
brew tap k8scat/tap
# 安装
brew install acli

# 一条命令直接安装
brew install k8scat/tap/acli

# 后续升级
brew update
brew upgrade k8scat/tap/acli
```

### Docker

```shell
# 将配置文件的目录挂载到容器内
docker run \
  -it \
  --rm \
  -v $HOME/.config/articli:/root/.config/articli \
  k8scat/articli:latest \
  juejin auth login

# 升级
docker pull k8scat/articli:latest
```

### 二进制

Please download from the [releases page](https://github.com/k8scat/Articli/releases).

### 源码编译

```shell
git clone https://github.com/k8scat/articli.git
cd articli
make
```

## 文档

[Docs](./docs/README.md)

## LICENSE

[MIT](./LICENSE)
