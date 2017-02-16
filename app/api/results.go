package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Results is a generic wrapper for all returned objects from the markup
// APIish endpoints.
type Results struct {
	Err     error       `json:"error,omitempty"`
	Results interface{} `json:"results"`
}

// RenderTo writes the results to a http.ResponseWriter as a JSON string.
// If Err is set it will render the error and return HTTP responde code 500.
func (r Results) RenderTo(rw http.ResponseWriter) {
	if r.Err != nil {
		txt := fmt.Sprintf("could not process request: %s", r.Err.Error())
		errMsg := fmt.Sprintf("{\n  \"error\": %q\n}", txt)
		http.Error(rw, errMsg, http.StatusInternalServerError)
		return
	}

	b, err := json.MarshalIndent(r.Results, "", "  ")
	if err != nil {
		errMsg := fmt.Sprintf("{\n  \"error\": \"could not process search: %q\"\n}", err.Error())
		http.Error(rw, errMsg, http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(b)
}
