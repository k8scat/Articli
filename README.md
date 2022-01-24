# Articli

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

## Support

- [GitHub](https://github.com)
  - [x] 认证
    - [x] 登录
    - [x] 登出
    - [x] 查看状态
  - [x] 仓库文件
    - [x] 上传
    - [x] 列取
    - [x] 删除
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

### NPM

```shell
npm install -g @k8scat/articli
```

### Homebrew

```shell
# 使用 tap
brew tap k8scat/tap
# 安装 Articli
brew install acli

# 直接安装
brew install k8scat/tap/acli

# 升级
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

## 文章模板

我们将使用文件内容开头 `---` 之间的数据作为文章的配置信息（元数据），
根据配置信息在不同平台上创建或更新文章，参考 [文章模板](./templates/article.md)。

## 使用说明

所有的命令都可以通过 `-h` 或 `--help` 参数查看帮助信息。

### GitHub

#### 登录

使用 GitHub Token 进行登录

```shell
# 交互式登录
acli github auth login

# 从标准输入获取 Token
acli github auth login --with-token < token.txt
```

#### 上传文件

```shell
# 上传 README.md 文件到 testrepo 仓库
acli github file upload --repo testrepo README.md
```

#### 列取文件

```shell
# 获取代码仓 testrepo 根目录的文件列表，包括文件和目录
acli github file get --repo testrepo

# 如果 testpath 是目录，则获取代码仓 testrepo 中 testpath 目录下的文件；
# 如果 testpath 是文件，则只获取该文件
acli github file get --repo testrepo --path testpath
```

![articli-github-file-upload.png](https://raw.githubusercontent.com/storimg/img/master/k8scat.com/articli-github-file-get.png)

#### 删除文件

```shell
# 使用 -o 或 --owner 可以指定仓库的 owner
acli github file delete --owner testowner --repo testrepo --path testdir/filename.txt
```

### 掘金 CLI

#### 登录

使用浏览器 Cookie 进行登录

```shell
# 交互式登录
acli juejin auth login

# 从标准输入获取 Cookie
acli juejin auth login --with-cookie < cookie.txt
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
acli juejin category list
```

#### 查看标签

```shell
# 过滤关键字
acli juejin tag list -k Go
```

#### 缓存标签

由于标签的数量比较多，可以通过设置缓存加快读取速度

```shell
# 设置缓存
acli juejin tag cache

# 使用缓存
acli jujin tag list --use-cache
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
alias gcli="acli github"
EOF

# 生效
source ~/.bashrc

# 使用简化后的命令查看掘金的登录状态
jcli auth status
```

## LICENSE

[MIT](./LICENSE)
