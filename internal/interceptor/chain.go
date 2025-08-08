package interceptor

import (
	"google.golang.org/grpc"
)

type InterceptorChain struct {
	interceptors []UnaryInterceptor
}

func NewInterceptorChain() *InterceptorChain {
	return &InterceptorChain{
		interceptors: make([]UnaryInterceptor, 0),
	}
}

func (ic *InterceptorChain) Add(interceptor UnaryInterceptor) *InterceptorChain {
	ic.interceptors = append(ic.interceptors, interceptor)
	return ic
}

func (ic *InterceptorChain) Build() []grpc.UnaryServerInterceptor {
	result := make([]grpc.UnaryServerInterceptor, len(ic.interceptors))
	for i, interceptor := range ic.interceptors {
		result[i] = interceptor.Intercept()
	}
	return result
}
