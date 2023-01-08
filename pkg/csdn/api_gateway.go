package csdn

import (
	"github.com/juju/errors"
	sign "github.com/k8scat/aliyun-api-gateway-sign-golang"
)

var (
	ResourceGateway *sign.APIGateway
	UserGateway     *sign.APIGateway
)

func initResourceGateway() error {
	gateway, err := sign.NewAPIGateway(ResourceAppKey, ResourceAppSecret)
	if err != nil {
		return errors.Trace(err)
	}
	ResourceGateway = gateway
	return nil
}

func initUserGateway() error {
	gateway, err := sign.NewAPIGateway(UserAppKey, UserAppSecret)
	if err != nil {
		return errors.Trace(err)
	}
	UserGateway = gateway
	return nil
}
