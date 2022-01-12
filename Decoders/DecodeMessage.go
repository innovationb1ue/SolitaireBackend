package Decoders

import (
	"encoding/json"
	"net/http"
)

// Msg2Map phrases byte stream as JSON string and return a map
func Msg2Map(msg []byte) map[string]interface{} {
	res := map[string]interface{}{}
	_ = json.Unmarshal(msg, &res)
	return res
}

// Req2Json phrases http.Response as JSON string and return a map
func Req2Json(r *http.Request) map[string]interface{} {
	ReqBodyJson := map[string]interface{}{}
	decode := json.NewDecoder(r.Body)
	_ = decode.Decode(&ReqBodyJson)
	return ReqBodyJson
}
