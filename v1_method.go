package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type eventPayload struct {
	IgnoreMessage bool `json:"ignore_message"`
}

type EventSubscription struct {
	Name          string
	IgnoreMessage bool
}

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

	if code != http.StatusOK && code != http.StatusCreated {
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

	if code != http.StatusCreated && code != http.StatusNoContent {
		return newUnexpectedResponseError(code, body)
	}

	return nil
}

func (this *V1) SubscribeDeviceEvents(deviceId string, events ...string) error {
	var empty interface{}

	payload := make(map[string]interface{})

	for _, v := range events {
		payload[v] = empty
	}

	body := newSubscribeDeviceEventsPayload(payload)

	path := "/subscriber/" + deviceId + "/subscriptions"
	code, res, postErr := this.request.post(path, "application/x-www-form-urlencoded", body)
	if postErr != nil {
		return postErr
	}

	if code == http.StatusBadRequest {
		return newUnexpectedResponseError(INVALID_PARAMETER_ERROR, res)
	}

	if code == http.StatusNotFound {
		return newUnexpectedResponseError(ACCOUNT_NOT_EXISTS_ERROR, res)
	}

	if code != http.StatusNoContent {
		return newUnexpectedResponseError(ACCOUNT_ALREADY_EXISTS_ERROR, res)
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

func (this *V1) UnsubscribeDeviceEvent(deviceId, event string) error {
	path := "/subscriber/" + deviceId + "/subscriptions/" + event
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

func (this *V1) GetDeviceEvents(deviceId string) ([]EventSubscription, error) {
	path := "/subscriber/" + deviceId + "/subscriptions"
	code, body, err := this.request.get(path)

	if err != nil {
		return nil, err
	}

	if code != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code: %d", code)
	}

	var raw map[string]eventPayload

	err = json.Unmarshal([]byte(body), &raw)

	if err != nil {
		return nil, err
	}

	var response []EventSubscription

	for k, v := range raw {

		event := EventSubscription{
			Name:          k,
			IgnoreMessage: v.IgnoreMessage,
		}

		response = append(response, event)
	}

	return response, nil
}

/**
 * Notifies all subscribers which are subscribed to the given language and sends them the given
 * message and data.
 */
func (this *V1) NotifyDevices(event, lang, msg, title string, data map[string]string, incrementBadge bool, category string, sound string, contentAvailable bool) error {
	localizedMsg := map[string]string{}
	localizedMsg[lang] = msg
	return this.NotifyDevicesLocalized(event, "", title, localizedMsg, data, incrementBadge, category, sound, contentAvailable)
}

/**
 * Sends a message to all subscribers.
 *
 * The message can be translated via the `localizedMsg` map. The order of lookups is as
 * following (assuming a subscription for `de-DE`):
 * 1. `localizedMap['de-DE']`
 * 2. `localizedMap['de']`
 * 3. `msg`
 *
 * If still no message is found, the subscriber is _not_ notified. To send a message only to subscribers
 * with a certain locale, leave the `msg` empty.
 */
func (this *V1) NotifyDevicesLocalized(event, msg, title string, localizedMsg, data map[string]string, incrementBadge bool, category string, sound string, contentAvailable bool) error {
	requestPayload := newNotifyPushNotificationRequestPayload(msg, title, localizedMsg, data, incrementBadge, category, sound, contentAvailable)
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
