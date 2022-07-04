package grpcserver

import (
	pb "sp-office-lookuper/pkg/protobuf"

	"github.com/asaskevich/govalidator"
)

type GetSortpointIDForm struct {
	DstOfficeID int64 `valid:"type(int64),required"`
}

func (f *GetSortpointIDForm) LoadAndValidate(req *pb.GetSortpointIdRequest) (valid bool, err error) {
	f.DstOfficeID = req.DstOfficeId

	_, err = govalidator.ValidateStruct(f)
	if err != nil {
		return false, err
	}

	return true, nil
}
