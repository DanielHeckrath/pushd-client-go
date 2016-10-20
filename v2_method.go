package client

import (
	context "golang.org/x/net/context"
	"net/http"
)

// Impact client version one.
type V2 struct {
	request *Request
}

type Badge struct {
	Value int64
}

type NotifyDevicesRequest struct {
	Event         string
	Message       string
	Title         string
	Localizations map[string]string
	Data          map[string]string

	IncrementBadge bool
	Badge          *Badge

	Category         string
	Sound            string
	ContentAvailable bool
}

func (v *V2) NotifyDevices(ctx context.Context, request NotifyDevicesRequest) error {
	requestPayload := newNotificationPayload(request)
	path := "/event/" + request.Event
	code, body, postErr := v.request.post(path, "application/x-www-form-urlencoded", requestPayload)
	if postErr != nil {
		return postErr
	}

	if code != http.StatusNoContent {
		return newUnexpectedResponseError(code, body)
	}

	return nil
}
