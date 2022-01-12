package Decoders

import (
	"encoding/json"
	"io"
)

// Msg2Map phrases byte stream as a JSON string and return a map
func Msg2Map(msg []byte) (map[string]interface{}, error) {
	res := map[string]interface{}{}
	_ = json.Unmarshal(msg, &res)
	return res, nil
}

// Req2Json phrases Body of http.Response as a JSON string and return a map
func Req2Json(body io.Reader) (map[string]interface{}, error) {
	ReqBodyJson := map[string]interface{}{}
	decode := json.NewDecoder(body)
	_ = decode.Decode(&ReqBodyJson)
	return ReqBodyJson, nil
}
