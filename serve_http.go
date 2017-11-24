package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func ServeResource(w http.ResponseWriter, r *http.Request, rc ResourceController, m Resource) {
	switch r.Method {
	case "GET":
		res, err := rc.ReadCollection(m)
		if err != nil {
			HTTPError(w, r, err)
		}

		err = HTTPSendJSON(w, res, http.StatusOK)
		if err != nil {
			log.Print("Couldn't convert your error to JSON.")
			return
		}
	default:
		// Unsupported Method, send Allow header back in error response
		w.Header().Add("Allow", "GET, PUT, DELETE")
		//msg := fmt.Sprintf("Method \"%s\" is not supported at %s.", r.Method, r.URL.Path)
		//HTTPError(w, r, &ApplicationError{Msg: msg, Code: http.StatusMethodNotAllowed})
		// todo: application error with code?
		HTTPError(w, r, nil)
		return
	}
}

// HTTPSendResource will marshal a single given Resource, set the correct status code and send the response.
func HTTPSendJSON(w http.ResponseWriter, m interface{}, code int) error {
	// MarshalIndent makes sure our JSON is pretty
	res, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}

	// But MarshalIndent doesn't come with a newline, so we do that ourselves.
	res = append(res, "\n"...)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	// Returned int is number of bytes written, not needed, so ignore
	_, err = w.Write(res)
	if err != nil {
		return err
	}

	return nil
}

// HTTPError returns an application-formatted JSON error object for when bad things happen.
func HTTPError(w http.ResponseWriter, r *http.Request, e error) {
	log.Printf("Error: %s\n", error.Error)

	// Marshal the error string to send in our response
	je, jerr := json.Marshal(e)
	if jerr != nil {
		log.Println("Couldn't convert your error to JSON.")
		return
	}

	w.WriteHeader(400)
	w.Header().Set("Content-Type", "application/json")
	w.Write(je)
}
