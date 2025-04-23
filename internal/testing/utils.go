package testing

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
)

func UnmarshalBodyArray[S ~[]E, E comparable](responseRecorder *httptest.ResponseRecorder) S {
	var unmarshalArray S
	_ = json.Unmarshal(responseRecorder.Body.Bytes(), &unmarshalArray)

	return unmarshalArray
}

func UnmarshalBody[E any](responseRecorder *httptest.ResponseRecorder) E {
	var unmarshalBody E
	_ = json.Unmarshal(responseRecorder.Body.Bytes(), &unmarshalBody)

	return unmarshalBody
}

func MarshalBody[E any](game E) io.Reader {
	marshalled, _ := json.Marshal(game)

	return bytes.NewReader(marshalled)
}
