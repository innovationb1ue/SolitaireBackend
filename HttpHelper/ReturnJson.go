package HttpHelper

import (
	"encoding/json"
	"net/http"
)

func ReturnJson(w http.ResponseWriter, v map[string]interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(v)
	return err
}
