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
	// Name Platform name
	Name() string
	// Auth Authenticate with raw auth data, like cookie or user:pass
	Auth(raw string) (username string, err error)
	// Publish Post article
	Publish(r io.Reader) (url string, err error)
	// ParseMark Parse markdown meta data
	ParseMark(mark *markdown.Mark) (params map[string]any, err error)
}
```

然后将新的平台注册到全局 `pltformHub` 中：

```go
// pkg/platform/hub.go
func init() {
	register(new(another.Platform))
}
```

## LICENSE

[MIT](./LICENSE)
