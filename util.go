package client

import (
	"bytes"
	"encoding/json"
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

func newNotifyPushNotificationRequestPayload(msg string, localizedMsg, data map[string]string) io.Reader {
	d := url.Values{}
	if msg != "" {
		d.Set("msg", msg)
	}
	for lang, localizedMsg := range localizedMsg {
		d.Set("msg."+lang, localizedMsg)
	}
	d.Set("sound", "default")

	for key, value := range data {
		d.Set("data."+key, value)
	}

	return bytes.NewBufferString(d.Encode())
}
