package util

import (
	"encoding/json"
	"net/http"
)

func PhraseJson(r *http.Request) map[string]interface{} {
	ReqBodyJson := map[string]interface{}{}
	decode := json.NewDecoder(r.Body)
	_ = decode.Decode(&ReqBodyJson)
	return ReqBodyJson
}
