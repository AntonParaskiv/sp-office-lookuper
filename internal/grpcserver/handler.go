package grpcserver

import (
	"sp-office-lookuper/internal/logging"
	pb "sp-office-lookuper/pkg/protobuf"
)

type grpcServerHandle struct {
	storage Storage
	logger  *logging.Logger
	pb.UnimplementedTransferBoxApiServer
}

func NewGRPCHandler(storage Storage, logger *logging.Logger) pb.TransferBoxApiServer {
	return &grpcServerHandle{
		storage:                           storage,
		logger:                            logger,
		UnimplementedTransferBoxApiServer: pb.UnimplementedTransferBoxApiServer{},
	}
}
