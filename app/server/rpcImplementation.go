package server

import (
	"Config/app/proto"
	"context"
	"fmt"
)

type ConfigWrapper struct {
	proto.UnimplementedConfigWrapperServer
}

func (receiver *ConfigWrapper) Get(ctx context.Context, in *proto.Service) (
	*proto.Data,
	error,
) {
	fmt.Println("Get")
	fmt.Println(in.GetService())
	return nil, nil
}
func (receiver *ConfigWrapper) Create(ctx context.Context, in *proto.Config) (*proto.Service, error) {
	fmt.Println("Create")

	return nil, nil
}

func (receiver *ConfigWrapper) Delete(context.Context, *proto.Service) (*proto.Service, error) {
	return nil, nil
}
