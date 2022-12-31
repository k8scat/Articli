# Articli

[![GitHub Repo stars](https://img.shields.io/github/stars/k8scat/articli?style=social)](https://github.com/k8scat/Articli/stargazers)
[![GitHub watchers](https://img.shields.io/github/watchers/k8scat/articli?style=social)](https://github.com/k8scat/Articli/watchers)
[![codecov](https://codecov.io/gh/k8scat/Articli/branch/main/graph/badge.svg?token=045FCRVF27)](https://codecov.io/gh/k8scat/Articli)

**Articli** 通过解析 `Markdown` 文件内容以及调用不同平台的接口，实现内容快速在不同平台进行发布。

## 平台

- [掘金](https://juejin.cn)
- [CSDN](https://csdn.net)
- [开源中国](https://oschina.net)
- [思否](https://segmentfault.com)

## 文档

https://k8scat.github.io/Articli

## 开发

如果您想添加其他平台，其实很简单，只需实现以下接口即可：

```go
type Platform interface {
    Name() string
    Auth(raw string) (string, error)
    NewArticle(r io.Reader) error
    Publish() (string, error)
}
```

然后将新的平台注册到全局 `platformHub` 中：

```go
// pkg/platform/hub.go
func init() {
    register(new(another.Platform))
}
```

## LICENSE

[MIT](./LICENSE)
