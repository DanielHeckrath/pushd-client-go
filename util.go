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

func newSubscribeDeviceRequestPayload(proto, token, lang string) io.Reader {
	data := url.Values{}
	data.Set("proto", proto)
	data.Set("token", token)
	data.Set("lang", lang)

	return bytes.NewBufferString(data.Encode())
}

func newNotifyPushNotificationRequestPayload(lang, msg string, data map[string]string) io.Reader {
	d := url.Values{}
	d.Set("msg."+lang, msg)

	for key, value := range data {
		d.Set("data."+key, value)
	}

	return bytes.NewBufferString(d.Encode())
}
