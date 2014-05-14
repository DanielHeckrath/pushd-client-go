package client

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
)

func createJsonPayload(v interface{}) (io.Reader, error) {
	payloadBytes, err := json.Marshal(v)
	if err != nil {
		return strings.NewReader(""), err
	}

	return strings.NewReader(string(payloadBytes)), nil
}

func newSubscribeDeviceRequestPayload(provider, token, language string) io.Reader {
	data := url.Values{}
	data.Set("proto", provider)
	data.Set("token", token)
	data.Set("lang", language)

	return bytes.NewBufferString(data.Encode())
}

func newFacebookInviteRequestPayload(userId, invitedFbId string) (io.Reader, error) {
	requestPayload := map[string]string{
		"user_id":             userId,
		"invited_facebook_id": invitedFbId,
	}

	return createPayload(requestPayload)
}

func newFacebookConnectRequestPayload(userId, userFbId string) (io.Reader, error) {
	requestPayload := map[string]string{
		"user_id":          userId,
		"user_facebook_id": userFbId,
	}

	return createPayload(requestPayload)
}

func newGetFriendListsRequestPayload(userId string, depth int) (io.Reader, error) {
	requestPayload := map[string]interface{}{
		"user_id": userId,
		"depth":   depth,
	}

	return createPayload(requestPayload)
}
