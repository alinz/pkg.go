package json

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	emptyBody = []byte("")
)

type PrepareResponser interface {
	PrepareResponse()
}

func Response(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if payload == nil || status == http.StatusNoContent {
		w.Write(emptyBody)
		return
	}

	if err, ok := payload.(error); ok {
		errorPayload := struct {
			Error string `json:"error"`
		}{}

		if jsonErr, ok := err.(*json.UnmarshalTypeError); ok {
			errorPayload.Error = fmt.Sprintf("field %s has wrong value", jsonErr.Field)
		} else {
			errorPayload.Error = err.Error()
		}

		payload = errorPayload
	} else if p, ok := payload.(PrepareResponser); ok {
		p.PrepareResponse()
	}

	Writer(w, payload)
}
