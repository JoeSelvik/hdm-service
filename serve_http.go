package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func ServeResource(w http.ResponseWriter, r *http.Request, rc ResourceController, m Resource) {
	switch r.Method {
	case "GET":
		res, aerr := rc.ReadCollection()
		if aerr != nil {
			HTTPError(w, r, aerr)
		}

		aerr = HTTPSendJSON(w, res, http.StatusOK)
		if aerr != nil {
			// todo: How to gracefully handle and log server side errors
			log.Println("Couldn't write back to the client.")
			return
		}
	default:
		// Unsupported Method, send Allow header back in error response
		w.Header().Add("Allow", "GET, PUT, DELETE")
		msg := fmt.Sprintf("Method \"%s\" is not supported at %s.", r.Method, r.URL.Path)
		HTTPError(w, r, &ApplicationError{Msg: msg, Code: http.StatusMethodNotAllowed})
		return
	}
}

// HTTPSendResource marshals given Resources, sets the correct status code, and sends the response.
func HTTPSendJSON(w http.ResponseWriter, m interface{}, code int) *ApplicationError {
	// MarshalIndent makes sure our JSON is pretty
	res, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		msg := "Couldn't format object into JSON"
		return &ApplicationError{Msg: msg, Err: err, Code: http.StatusUnsupportedMediaType}
	}

	// But MarshalIndent doesn't come with a newline, so we do that ourselves.
	res = append(res, "\n"...)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	// Returned int is number of bytes written, not needed, so ignore
	_, err = w.Write(res)
	if err != nil {
		msg := "Couldn't write back to client."
		return &ApplicationError{Msg: msg, Err: err, Code: http.StatusUnsupportedMediaType}
	}

	return nil
}

// HTTPError returns an application-formatted JSON error object for when bad things happen.
func HTTPError(w http.ResponseWriter, r *http.Request, aerr *ApplicationError) {
	log.Printf("Error: %s\n", aerr.Error)

	// Marshal the error string to send in our response
	je, jerr := json.Marshal(aerr)
	if jerr != nil {
		// todo: how to handle? Panic? Print more info?
		log.Println("Couldn't convert your error to JSON.")
		return
	}

	w.WriteHeader(aerr.Code)
	w.Header().Set("Content-Type", "application/json")
	w.Write(je)
}
