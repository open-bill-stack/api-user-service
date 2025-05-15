package router

import (
	"google.golang.org/grpc"
)

type Router interface {
	Register(*grpc.Server)
}
