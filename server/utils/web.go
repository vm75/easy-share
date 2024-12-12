package utils

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func GetVariables(r *http.Request) map[string]string {
	return mux.Vars(r)
}

func GetVariable(r *http.Request, name string) string {
	return GetVariables(r)[name]
}

func GetParams(r *http.Request) map[string]string {
	params := make(map[string]string)
	for k, v := range r.URL.Query() {
		if len(v) == 0 {
			continue
		}
		params[k] = v[0]
	}
	return params
}

func GetParam(r *http.Request, name string) string {
	return GetParams(r)[name]
}

func GetJsonContent(r *http.Request) (map[string]interface{}, error) {
	var content map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&content)
	return content, err
}

func GetContent(r *http.Request, content interface{}) error {
	return json.NewDecoder(r.Body).Decode(&content)
}
