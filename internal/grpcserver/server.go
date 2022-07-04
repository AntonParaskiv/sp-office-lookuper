package grpcserver

import (
	"context"
	"math"
	"net"

	"sp-office-lookuper/internal/logging"
	pb "sp-office-lookuper/pkg/protobuf"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	DefaultServerMaxReceiveMessageSize = math.MaxInt32
	DefaultServerMaxSendMessageSize    = 1024 * 1024 * 24
)

type Storage interface {
	SetOfficeSortPoint(officeID, sortPointID int64)
	GetSortPoint(officeID int64) (int64, error)
}

type Server struct {
	grpcServer *grpc.Server
	handlers   pb.TransferBoxApiServer
	listener   net.Listener
	storage    Storage
	logger     *logging.Logger
}

func (s *Server) Close(ctx context.Context) {
	s.logger.Info(ctx, "GRPC server shutting down...")
	s.grpcServer.Stop()
}

func (s *Server) ListenAndServe() (err error) {
	options := DefaultServeOptions()
	msgSizeOptions := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(DefaultServerMaxReceiveMessageSize),
		grpc.MaxSendMsgSize(DefaultServerMaxSendMessageSize),
	}
	options = append(options, msgSizeOptions...)
	options = append(options, grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer())))

	s.grpcServer = grpc.NewServer(options...)
	reflection.Register(s.grpcServer)
	pb.RegisterTransferBoxApiServer(s.grpcServer, s.handlers)

	return s.grpcServer.Serve(s.listener)
}

func NewServer(listener net.Listener, storage Storage, logger *logging.Logger) (Server, error) {
	srv := Server{
		handlers: NewGRPCHandler(storage, logger),
		listener: listener,
		storage:  storage,
		logger:   logger,
	}

	return srv, nil
}
