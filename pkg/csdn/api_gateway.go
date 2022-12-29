package csdn

import (
	sign "github.com/k8scat/aliyun-api-gateway-sign-golang"
)

var (
	ResourceGateway *sign.APIGateway
	UserGateway     *sign.APIGateway
)

func InitResourceGateway() error {
	gateway, err := sign.NewAPIGateway(ResourceAppKey, ResourceAppSecret)
	if err != nil {
		return err
	}
	ResourceGateway = gateway
	return nil
}

func InitUserGateway() error {
	gateway, err := sign.NewAPIGateway(UserAppKey, UserAppSecret)
	if err != nil {
		return err
	}
	UserGateway = gateway
	return nil
}
