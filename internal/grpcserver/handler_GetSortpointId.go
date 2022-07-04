package grpcserver

import (
	"context"
	"errors"

	"sp-office-lookuper/internal/app"
	"sp-office-lookuper/internal/tracer"
	pb "sp-office-lookuper/pkg/protobuf"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	OperationGetSortpointID = "GetSortpointId"
)

//nolint:revive,stylecheck // it depends on proto-file
func (h *grpcServerHandle) GetSortpointId(ctx context.Context,
	req *pb.GetSortpointIdRequest) (*pb.GetSortpointIdResponse, error) {
	var err error
	ctx, span, finish := tracer.Trace(ctx)
	defer func() { finish(err) }()
	span.SetTag(app.OperationTag, OperationGetSortpointID)
	span.SetTag(app.DstOfficeIDTag, req.DstOfficeId)

	logEntry := h.logger.CreateEntry().
		WithField(app.OperationTag, OperationGetSortpointID).
		WithField(app.DstOfficeIDTag, req.DstOfficeId)

	form := GetSortpointIDForm{}
	valid, err := form.LoadAndValidate(req)
	if err != nil {
		logEntry.WithError(err).Error(ctx, "unable load and validate request values")
		if !valid {
			return nil, status.Error(codes.InvalidArgument, "bad request")
		}
		return nil, status.Error(codes.Internal, "something went wrong")
	}

	sortPointID, err := h.storage.GetSortPoint(form.DstOfficeID)
	if err != nil {
		if errors.Is(err, app.ErrOfficeNotFound) {
			logEntry.Warn(ctx, "office not found")
			return nil, err
		}
		logEntry.WithError(err).Error(ctx, "get sort point from storage failed")
		return nil, err
	}

	span.SetTag(app.SortPointIDTag, sortPointID)
	logEntry = logEntry.
		WithField(app.SortPointIDTag, sortPointID)

	resp := &pb.GetSortpointIdResponse{
		SortpointId: sortPointID,
	}

	logEntry.Info(ctx, "successfully handled")
	return resp, nil
}
