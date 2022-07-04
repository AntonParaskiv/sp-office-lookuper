package grpcserver

import (
	"time"

	"google.golang.org/grpc"
	grpcKeepalive "google.golang.org/grpc/keepalive"
)

func DefaultServeOptions() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.KeepaliveEnforcementPolicy(DefaultEnforcementServerOptions()),
		grpc.KeepaliveParams(DefaultKeepaliveServerOptions()),
	}
}

func DefaultEnforcementServerOptions() grpcKeepalive.EnforcementPolicy {
	return grpcKeepalive.EnforcementPolicy{
		MinTime:             5 * time.Second,
		PermitWithoutStream: true,
	}
}

func DefaultKeepaliveServerOptions() grpcKeepalive.ServerParameters {
	return grpcKeepalive.ServerParameters{
		MaxConnectionIdle:     15 * time.Second,
		MaxConnectionAge:      30 * time.Second,
		MaxConnectionAgeGrace: 5 * time.Second,
		Time:                  5 * time.Second,
		Timeout:               1 * time.Second,
	}
}
