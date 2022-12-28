package csdn

import "github.com/juju/errors"

// csdn图床
// https://blog.51cto.com/144dotone/2952199
//
// 关于CSDN获取博客内容接口的x-ca-signature签名算法研究
// https://blog.csdn.net/chouzhou9701/article/details/109306318
//
// 阿里云 - 使用摘要签名认证方式调用API
// https://help.aliyun.com/document_detail/29475.html
//
// https://github.com/aliyun/api-gateway-demo-sign-python
//
// https://github.com/k8scat/aliyun-api-gateway-sign-golang

const (
	// https://csdnimg.cn/release/md/static/js/app.chunk.463d2f5b.js
	ResourceAppKey    = "203803574"
	ResourceAppSecret = "9znpamsyl2c7cdrr9sas0le9vbc3r6ba"

	// https://i.csdn.net/static/js/app.c55be3.js
	UserAppKey    = "203796071"
	UserAppSecret = "i5rbx2z2ivnxzidzpfc0z021imsp2nec"

	BizAPIBase = "https://bizapi.csdn.net"
)

type Client struct {
	cookie string
}

func New(cookie string) (*Client, error) {
	if cookie == "" {
		return nil, errors.New("Invalid cookie")
	}
	return &Client{
		cookie: cookie,
	}, nil
}
