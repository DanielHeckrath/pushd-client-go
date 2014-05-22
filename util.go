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

func newSubscribeDeviceRequestPayload(proto, token, lang string) io.Reader {
	data := url.Values{}
	data.Set("proto", proto)
	data.Set("token", token)
	data.Set("lang", lang)
	data.Set("badge", "0")

	return bytes.NewBufferString(data.Encode())
}

func newNotifyPushNotificationRequestPayload(lang, msg string, data map[string]string) io.Reader {
	d := url.Values{}
	d.Set("msg."+lang, msg)
	d.Set("sound", "default")

	for key, value := range data {
		d.Set("data."+key, value)
	}

	fmt.Printf("pushd notification payload: %#v\n", d.Encode())

	return bytes.NewBufferString(d.Encode())
}
