# Articli

[![Release](https://github.com/k8scat/Articli/actions/workflows/release.yaml/badge.svg)](https://github.com/k8scat/Articli/actions/workflows/release.yaml)
[![GitHub Repo stars](https://img.shields.io/github/stars/k8scat/articli?style=social)](https://github.com/k8scat/Articli/stargazers)
[![GitHub watchers](https://img.shields.io/github/watchers/k8scat/articli?style=social)](https://github.com/k8scat/Articli/watchers)
[![codecov](https://codecov.io/gh/k8scat/Articli/branch/main/graph/badge.svg?token=045FCRVF27)](https://codecov.io/gh/k8scat/Articli)

**Articli** 通过解析 `Markdown` 文件内容以及调用不同平台的接口，实现内容快速在不同平台进行发布。

## 支持的平台

### 平台文章管理

- [掘金](https://juejin.cn)
- [CSDN](https://csdn.net)

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

## 文章模板

我们将使用文件内容开头 `---` 之间的数据作为文章的配置信息（元数据），
根据配置信息在不同平台上创建或更新文章，参考 [文章模板](https://raw.githubusercontent.com/k8scat/Articli/csdn/templates/article.md)。

```markdown
---
# 通用配置，其他平台可以继承该配置
title: 标题1
brief_content: 内容概要
cover_image:
- https://img.alicdn.com/tfs/TB1.jpg
prefix_content: ""  # 前缀内容
suffix_content: |
  ## Powered by

  本文由 [Articli](https://github.com/k8scat/Articli.git) 工具自动发布。

juejin:
  article_id: ""  # 手动填写
  title: 标题2 # 如果不填写，则使用通用配置中的 title
  tag_ids:
  - "6809640675944955918"
  category_id: "6809637771511070734"
  cover_image: https://img.alicdn.com/tfs/TB1.jpg
  brief_content: 内容概要
  prefix_content: "这是我参与xx活动..." # 前缀内容，主要用于掘金的活动
  suffix_content: |
    ## Powered by

    本文由 [Articli](https://github.com/k8scat/Articli.git) 工具自动发布。
  sync_to_org: false # 是否同步到组织，个人账号不支持

csdn:
  article_id: ""  # 手动填写
  title: 标题3
  brief_content: 内容概要
  categories:
  - Golang
  - 后端
  tags:
  - cli
  - csdn
  # 可选值: public, private, read_need_vip, read_need_fans
  read_type: public
  # 可选值: 发布 publish, 草稿 draft
  publish_status: publish
  # 可选值: 原创 original, 转载 repost, 翻译 translated
  article_type: original
  # 转载时必须填写
  original_url: ""
  # 原文允许转载或者本次转载已经获得原文作者授权
  authorized_status: false
  # 支持单图、三图、无图
  cover_images:
  - https://img.alicdn.com/tfs/TB1.jpg
  - https://img.alicdn.com/tfs/TB2.jpg
  - https://img.alicdn.com/tfs/TB3.jpg
  prefix_content: "这是我参与xx活动..." # 前缀内容，主要用于掘金的活动
  suffix_content: |
    ## Powered by

    本文由 [Articli](https://github.com/k8scat/Articli.git) 工具自动发布。
---

内容概要

<!-- more -->

正文内容
```

## 使用说明

所有的命令都可以通过 `-h` 或 `--help` 参数查看帮助信息。

```shell
$ acli --help
Publish article anywhere.

Usage:
  acli [command]

Available Commands:
  auth        Authenticate
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  pub         Publish article
  version     Show version information

Flags:
  -c, --config string     Config file
  -h, --help              help for acli
  -p, --platform string   Platform name

Use "acli [command] --help" for more information about a command.
```

### 查看版本

```shell
acli version
```


### 登录账号

使用浏览器 Cookie 进行登录

```shell
# platform: juejin, csdn
acli auth -p <platform> --raw <cookie>
```

### 发布文章

```shell
acli pub -p <platform> --file <article-file>
```
