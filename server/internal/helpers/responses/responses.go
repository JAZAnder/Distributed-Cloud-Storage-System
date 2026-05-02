package responses

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/JAZAnder/Distributed-Cloud-Storage-System/server/internal/helpers/quickLog"

)

func RespondWithError(r *http.Request, w http.ResponseWriter, code int, message string) {
	quickLog.Log("responseMaster", "HTTP Response", "", r.RemoteAddr, "", strconv.Itoa(code), message)
	RespondWithJSONNoLog(w, code, map[string]string{"error": message})
}

func RespondWithJSON(r *http.Request, w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
	quickLog.Log("responseMaster", "HTTP Response", "", r.RemoteAddr, "", strconv.Itoa(code), string(response))
}

func RespondWithJSONNoLog(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
