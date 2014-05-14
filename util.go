package client

import (
	"bytes"
	"io"
	"net/url"
)

func newSubscribeDeviceRequestPayload(provider, token, language string) io.Reader {
	data := url.Values{}
	data.Set("proto", provider)
	data.Set("token", token)
	data.Set("lang", language)

	return bytes.NewBufferString(data.Encode())
}
