package httpserver

import (
	"io"
	"net/http"

	"sp-office-lookuper/internal/app"
	"sp-office-lookuper/internal/tracer"
)

type Storage interface {
	SetOfficeSortPoint(officeID, sortPointID int64)
	GetSortPoint(officeID int64) (int64, error)
}

func (h *HTTPHandlers) OfficeHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx, span, finish := tracer.Trace(r.Context(), tracer.ExtractFrom(r.Header))
	defer func() { finish(err) }()
	span.SetTag(app.OperationTag, OperationOffice)

	logEntry := h.logger.CreateEntry().
		WithField(app.OperationTag, OperationOffice)

	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logEntry.
			WithError(err).
			Error(ctx, "failed to validate form")
	}

	form := &OfficeForm{}
	if err = form.LoadAndValidate(body); err != nil {
		span.SetTag(app.BodyTag, string(body))
		logEntry.
			WithField(app.BodyTag, string(body)).
			WithError(err).
			Error(ctx, "failed to validate form")
		apiError(w, err.Error(), CodeInvalidArgument, http.StatusBadRequest)
		return
	}

	span.SetTag(app.SortPointIDTag, form.SortPointID)
	span.SetTag(app.DstOfficeIDTag, form.DstOfficeID)

	logEntry = logEntry.
		WithField(app.SortPointIDTag, form.SortPointID).
		WithField(app.DstOfficeIDTag, form.DstOfficeID)

	h.storage.SetOfficeSortPoint(form.DstOfficeID, form.SortPointID)

	logEntry.Info(ctx, "successfully handled")
}
