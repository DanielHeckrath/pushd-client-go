package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
)

func stringify(v interface{}) (string, error) {
	payloadBytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	return string(payloadBytes), nil
}

func newSubscribeDeviceEventsPayload(events map[string]interface{}) io.Reader {
	payload, err := json.Marshal(events)

	if err != nil {
		return &bytes.Buffer{}
	}

	return bytes.NewBuffer(payload)
}

func newSubscribeDeviceRequestPayload(proto, token, lang string) io.Reader {
	data := url.Values{}
	data.Set("proto", proto)
	data.Set("token", token)
	data.Set("lang", lang)
	data.Set("badge", "0")

	return bytes.NewBufferString(data.Encode())
}

func newNotifyPushNotificationRequestPayload(msg, title string, localizedMsg, data map[string]string, incrementBadge bool, category string, sound string, contentAvailable bool) io.Reader {
	d := url.Values{}
	if msg != "" {
		d.Set("msg", msg)
	}
	if title != "" {
		d.Set("title", title)
	}
	if !incrementBadge {
		d.Set("incrementBadge", "false")
	}
	for lang, localizedMsg := range localizedMsg {
		d.Set("msg."+lang, localizedMsg)
	}
	if sound == "" {
		sound = "default"
	}
	d.Set("sound", sound)
	if category != "" {
		d.Set("category", category)
	}

	if contentAvailable {
		d.Set("contentAvailable", "true")
	}

	for key, value := range data {
		d.Set("data."+key, value)
	}

	return bytes.NewBufferString(d.Encode())
}

func newNotificationPayload(request NotifyDevicesRequest) io.Reader {
	if request.ContentAvailable {
		var badge int64

		if request.Badge != nil {
			badge = request.Badge.Value
		}

		return newContentAvailablePayload(request.Data, badge)
	}

	d := url.Values{}
	if request.Message != "" {
		d.Set("msg", request.Message)
	}
	if request.Title != "" {
		d.Set("title", request.Title)
	}

	if request.Badge != nil {
		d.Set("badge", fmt.Sprintf("%d", request.Badge.Value))
		d.Set("incrementBadge", "false")
	} else {
		if !request.IncrementBadge {
			d.Set("incrementBadge", "false")
		}
	}

	for lang, localizedMsg := range request.Localizations {
		d.Set("msg."+lang, localizedMsg)
	}
	if request.Sound == "" {
		request.Sound = "default"
	}
	d.Set("sound", request.Sound)
	if request.Category != "" {
		d.Set("category", request.Category)
	}

	for key, value := range request.Data {
		d.Set("data."+key, value)
	}

	return bytes.NewBufferString(d.Encode())
}

func newContentAvailablePayload(data map[string]string, badge int64) io.Reader {
	d := url.Values{}

	d.Set("contentAvailable", "true")
	d.Set("incrementBadge", "false")
	d.Set("badge", fmt.Sprintf("%d", badge))

	for key, value := range data {
		d.Set("data."+key, value)
	}

	return bytes.NewBufferString(d.Encode())
}
