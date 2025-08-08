package server

import (
	"fmt"
	"net"

	proto "github.com/zhenyanesterkova/metricsmonitor/internal/app/proto/metric"
	"github.com/zhenyanesterkova/metricsmonitor/internal/grpchandler"
	"github.com/zhenyanesterkova/metricsmonitor/internal/interceptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func (s *Server) runGRPC() error {
	listen, err := net.Listen("tcp", s.config.SConfig.Address)
	if err != nil {
		s.logger.LogrusLog.Errorf("failed listen announces on the local network address: %v", err)
		return fmt.Errorf("failed listen announces on the local network address: %w", err)
	}

	loggerAdapter := interceptor.NewLogrusAdapter(s.logger)
	loggingInterceptor := interceptor.NewInterceptorLogger(loggerAdapter)

	var signatureInterceptor *interceptor.SignatureInterceptor
	if s.config.SConfig.HashKey != nil {
		validator := interceptor.NewHMACSignatureValidator(s.config.SConfig.HashKey, loggerAdapter)
		signatureInterceptor = interceptor.NewSignatureInterceptor(validator)
	}

	compressionHandler := interceptor.NewGzipCompressionHandler(loggerAdapter)
	compressionInterceptor := interceptor.NewCompressionInterceptor(compressionHandler)

	chain := interceptor.NewInterceptorChain().Add(loggingInterceptor)
	if signatureInterceptor != nil {
		chain.Add(signatureInterceptor)
	}
	chain.Add(compressionInterceptor)

	creds, err := credentials.NewServerTLSFromFile(s.config.SConfig.CertPath, s.config.SConfig.KeyPath)
	if err != nil {
		s.logger.LogrusLog.Errorf(
			"failed constructs TLS credentials from the input certificate file and key file for server: %v",
			err,
		)
		return fmt.Errorf(
			"failed constructs TLS credentials from the input certificate file and key file for server: %w",
			err,
		)
	}

	grpcServer := grpc.NewServer(
		grpc.Creds(creds),
		grpc.ChainUnaryInterceptor(chain.Build()...),
	)

	grpcHandler := grpchandler.New(s.storage, s.logger)
	proto.RegisterMonitorServer(grpcServer, grpcHandler)

	s.logger.LogrusLog.Infof("Start gRPC Server on %s", s.config.SConfig.Address)

	go func() {
		if err := grpcServer.Serve(listen); err != nil {
			s.logger.LogrusLog.Errorf("gRPC server error: %v", err)
		}
	}()

	return s.waitForShutdown(func() error {
		grpcServer.GracefulStop()
		return nil
	})
}
