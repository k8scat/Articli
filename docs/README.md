# Articli

[![GitHub Repo stars](https://img.shields.io/github/stars/k8scat/articli?style=social)](https://github.com/k8scat/Articli/stargazers)
[![GitHub watchers](https://img.shields.io/github/watchers/k8scat/articli?style=social)](https://github.com/k8scat/Articli/watchers)
[![codecov](https://codecov.io/gh/k8scat/Articli/branch/main/graph/badge.svg?token=045FCRVF27)](https://codecov.io/gh/k8scat/Articli)

**Articli** 通过解析 `Markdown` 文件内容以及调用不同平台的接口，实现内容快速在不同平台进行发布。

## 平台

- [掘金](https://juejin.cn)
- [CSDN](https://csdn.net)
- [开源中国](https://oschina.net)
- [SegmentFault](https://segmentfault.com)

## 安装

### 二进制

下载 [releases page](https://github.com/k8scat/Articli/releases).

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
# 通用配置，其他平台可以选择继承该配置，或者为不同平台进行单独设置
# 标题
title: "标题"
# 文章概要
brief_content: "文章概要"
# 封面图片地址
cover_images:
- https://img.alicdn.com/tfs/TB1.jpg
# 前缀内容
prefix_content: "文章前缀内容"
# 后缀内容
suffix_content: "文章后缀内容"

# 掘金平台文章配置
juejin:
    # 文章 id，不填写表示发布新文章
    article_id: ""
    # 草稿 id，不填则通过接口获取文章对应的草稿 id
    draft_id: ""
    # 标签名称
    tags:
    - "Go"
    - "Mac"
    # 分类名称
    category: "后端"
    # 是否同步到组织，个人账号不支持
    sync_to_org: false

csdn:
    # 文章 id，不填写表示发布新文章
    article_id: ""
    # 分类名称
    categories:
    - Golang
    - 后端
    # 标签名称
    tags:
    - cli
    - csdn
    # 发布形式，可选值：全部可见 public、仅我可见 private、VIP可见 read_need_vip、粉丝可见 read_need_fans，默认 public
    read_type: public
    # 发布状态，可选值：发布 publish、草稿 draft，默认 publish
    publish_status: publish
    # 文章类型，可选值：原创 original、转载 repost、翻译 translated，默认 original
    article_type: original
    # 原文链接，转载文章时必须填写
    original_url: ""
    # 原文允许转载或者本次转载已经获得原文作者授权
    authorized_status: false
    # 内容等级，可选择：初级 1、中级 2、高级 3，默认 1
    level: 1

oschina:
    # 文章 id，不填写表示发布新文章
    article_id: ""
    # 分类名称
    category: "Golang"
    # 推广专区名称
    technical_field_id: "程序人生"
    # 禁止评论
    deny_comment: false
    # 置顶
    top: false
    # 下载外站图片到本地
    download_image: false
    # 仅自己可见
    privacy: false

segmentfault:
    # 文章 id，不填写表示发布新文章
    article_id: ""
    # 标签名称
    tags:
    - "kubernetes"
    # 注明版权
    license: false
    # 文章类型，可选值：原创 1、转载 2、翻译 3，默认 1
    type: 1
    # 原文地址，如果是转载或翻译则必须填写
    url: ""

---

文章概要

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

### 登录账号

使用浏览器 cookie 进行登录：

- 掘金
- CSDN
- 开源中国

**思否请使用 token 进行登录（可以从浏览器请求头中获取）**

```shell
acli auth -p <platform> --raw <cookie>
```

### 发布文章

```shell
acli pub -p <platform> --file <article-file>
```
