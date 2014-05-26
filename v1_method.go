package client

import (
	"encoding/json"
	"net/http"
)

type SubscribeDeviceResponsePayload struct {
	Id      string `json:"id"`
	Created int    `json:"created"`
	Updated int    `json:"updated"`
	Proto   string `json:"proto"`
	Token   string `json:"token"`
	Lang    string `json:"lang"`
}

func (this *V1) SubscribeDevice(provider, token, language string) (SubscribeDeviceResponsePayload, error) {
	empty := SubscribeDeviceResponsePayload{}

	requestPayload := newSubscribeDeviceRequestPayload(provider, token, language)

	path := "/subscribers"
	code, body, postErr := this.request.post(path, "application/x-www-form-urlencoded", requestPayload)
	if postErr != nil {
		return empty, postErr
	}

	if code == http.StatusBadRequest {
		return empty, newUnexpectedResponseError(INVALID_PARAMETER_ERROR, body)
	}

	if code == http.StatusCreated {
		return empty, newUnexpectedResponseError(ACCOUNT_ALREADY_EXISTS_ERROR, body)
	}

	if code != http.StatusOK {
		return empty, newUnexpectedResponseError(code, body)
	}

	var responsePayload SubscribeDeviceResponsePayload
	if err := json.Unmarshal([]byte(body), &responsePayload); err != nil {
		return empty, err
	}

	return responsePayload, nil
}

func (this *V1) SubscribeDeviceEvent(deviceId, event string) error {
	path := "/subscriber/" + deviceId + "/subscriptions/" + event
	code, body, postErr := this.request.post(path, "application/x-www-form-urlencoded", nil)
	if postErr != nil {
		return postErr
	}

	if code == http.StatusBadRequest {
		return newUnexpectedResponseError(INVALID_PARAMETER_ERROR, body)
	}

	if code == http.StatusNotFound {
		return newUnexpectedResponseError(ACCOUNT_NOT_EXISTS_ERROR, body)
	}

	if code == http.StatusNoContent {
		return newUnexpectedResponseError(ACCOUNT_ALREADY_EXISTS_ERROR, body)
	}

	if code != http.StatusCreated {
		return newUnexpectedResponseError(code, body)
	}

	return nil
}

func (this *V1) UnsubscribeDevice(deviceId string) error {
	path := "/subscriber/" + deviceId
	code, body, postErr := this.request.del(path)
	if postErr != nil {
		return postErr
	}

	if code == http.StatusBadRequest {
		return newUnexpectedResponseError(INVALID_PARAMETER_ERROR, body)
	}

	if code == http.StatusNotFound {
		return newUnexpectedResponseError(ACCOUNT_NOT_EXISTS_ERROR, body)
	}

	if code != http.StatusNoContent {
		return newUnexpectedResponseError(code, body)
	}

	return nil
}

/**
 * Notifies all subscribers which are subscribed to the given language and sends them the given
 * message and data.
 */
func (this *V1) NotifyDevices(event, lang, msg string, data map[string]string) error {
	localizedMsg := map[string]string{}
	localizedMsg[lang] = msg
	return this.NotifyDevicesLocalized(event, "", localizedMsg, data)
}

/**
 * Sends a message to all subscribers.
 *
 * The message can be translated via the `localizedMsg` map. The order of lookups is as
 * following (assuming a subscription for `de_DE`):
 * 1. `localizedMap['de_DE']`
 * 2. `localizedMap['de']`
 * 3. `msg`
 *
 * If still no message is found, the subscriber is _not_ notified. To send a message only to subscribers
 * with a certain locale, leave the `msg` empty.
 */
func (this *V1) NotifyDevicesLocalized(event, msg string, localizedMsg, data map[string]string) error {
	requestPayload := newNotifyPushNotificationRequestPayload(msg, localizedMsg, data)
	path := "/event/" + event
	code, body, postErr := this.request.post(path, "application/x-www-form-urlencoded", requestPayload)
	if postErr != nil {
		return postErr
	}

	if code != http.StatusNoContent {
		return newUnexpectedResponseError(code, body)
	}

	return nil
}
