package httpserver

import (
	"github.com/asaskevich/govalidator"
)

type OfficeForm struct {
	SortPointID int64 `json:"sortpointId" valid:"type(int64),required"`
	DstOfficeID int64 `json:"dstOfficeId" valid:"type(int64),required"`
}

func (f *OfficeForm) LoadAndValidate(body []byte) error {
	err := f.UnmarshalJSON(body)
	if err != nil {
		return err
	}

	_, err = govalidator.ValidateStruct(f)
	if err != nil {
		return err
	}

	return nil
}
