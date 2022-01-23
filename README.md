# Articli

**Articli** is an Article CLI tool for managing content in multi platforms.

**Articli** 是一个可以管理多个平台内容的命令行工具，
通过解析 `Markdown` 文件内容以及调用平台接口，实现文章的发布、更新等功能，
目前仅支持[掘金](https://juejin.cn)，后续会继续支持其他平台。

最终目标是基于 **本地文件** + **Git 代码仓** 管理所有的文章，
并且可以通过命令行操作以及 CI/CD，实现文章在各个平台的发布、更新等功能。
这样做的好处有：

- 数据安全，既发布到了第三方平台，又可以通过 **Git 代码仓**管理，避免因平台问题导致数据丢失
- 可以实现自动化，比如文章自动在多个平台发布、更新
- 面向程序员的 CLI 工具，可以实现更多个性化的操作

为本项目点赞将鼓励作者继续完善下去，欢迎提出建议、Bug、PR。

## Support

- [掘金](https://juejin.cn)
  - [x] 认证
    - [x] 登录
    - [x] 登出
    - [x] 查看状态
  - [x] 文章
    - [x] 列取
    - [x] 删除
    - [x] 发布草稿
    - [x] 新建
    - [x] 更新
    - [x] 查看
  - [x] 图片
    - [x] 上传
  - [x] 草稿
    - [x] 创建
    - [x] 列取
    - [x] 删除
  - [x] 列取标签
  - [x] 列取分类
- [ ] [开源中国](https://oschina.net)
- [ ] [CSDN](https://csdn.net)

## 安装

Please download from the [releases page](https://github.com/k8scat/Articli/releases).

## 文章模板

我们将文件内容开头的 `---` 之间的数据作为文章的配置信息。

```markdown
---
# 通用配置，其他平台可以继承该配置
title: 标题1
brief_content: 内容概要
cover_image: https://img.alicdn.com/tfs/TB1.jpg

juejin:
  title: 标题2 # 如果不填写，则使用通用配置中的 title
  tags:
    - Go
    - 程序员
  category: 后端
  cover_image: https://img.alicdn.com/tfs/TB1.jpg
  brief_content: 内容概要
  prefix_content: "这是我参与xx活动..." # 前缀内容，主要用于掘金的活动
  suffix_content: |
    ## 原创申明
    
    本文由 `Articli` 工具自动发布。
  
  # 自动生成部分
  draft_id: "7xxx"
  draft_create_time: "2022-01-23 11:48:02"
  draft_update_time: "2022-01-24 11:48:02"
  article_id: "8xxx"
  article_create_time: "2022-01-25 11:48:02"
  article_update_time: "2022-01-26 11:48:02"

oschina:
  title: 标题3
  ...

csdn:
  title: 标题4
  ...
---

内容概要

<!-- more -->

正文内容
```

## 使用

所有的命令都可以通过 `-h` 或 `--help` 参数查看帮助信息。

### 掘金 CLI

#### 登录

使用浏览器 Cookie 进行登录

```shell
acli juejin auth login --with-cookie < cookie_file 
```

#### 创建/更新文章

```shell
# create 命令可以通过识别文章的配置信息，自动选择创建或者更新文章，同时发布到掘金
acli juejin article create markdown-file.md
```

#### 查看文章列表

通过 `-k` 或 `--keyword` 关键字参数过滤文章列表

```shell
acli juejin article list -k Docker
```

#### 打开文章

使用默认浏览器打开文章

```shell
acli juejin article view 7055689358657093646
```

#### 查看分类

```shell
acli juejin category list -k 后端
```

#### 查看标签

```shell
acli juejin tag list -k Go
```

#### 上传图片

支持上传本地图片和网络图片

```shell
# 本地图片
acli juejin image upload leetcode-go.png

# 网络图片
acli juejin image upload https://launchtoast.com/wp-content/uploads/2021/11/learn-rust-programming-language.png
```

### 简化命令

使用 `alias` 别名进行简化命令

```shell
# 将 acli juejin 简化成 jcli
cat >> ~/.bashrc << EOF
alias jcli="acli juejin"
EOF

# 生效
source ~/.bashrc

# 使用简化后的命令查看登录状态
jcli auth status
```

## LICENSE

[MIT](./LICENSE)
