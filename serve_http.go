package main

import (
	"encoding/json"
	"fmt"
	"github.com/JoeSelvik/hdm-service/models"
	"log"
	"net/http"
	"strconv"
)

func ServeResource(w http.ResponseWriter, r *http.Request, rc ResourceController, m models.Resource) {
	if len(r.URL.Path) > len(rc.Path()) {
		ServeSingleResource(w, r, rc, m)
	} else {
		ServeResourceCollection(w, r, rc, m)
	}
}

// ServeSingleResource redirects requests involving a specific resource via an id number.
func ServeSingleResource(w http.ResponseWriter, r *http.Request, rc ResourceController, m models.Resource) {
	// Parse id number
	id, err := strconv.Atoi(r.URL.Path[len(rc.Path()):])
	if err != nil {
		msg := "That doesn't look like a valid id number."
		HTTPError(w, r, &ApplicationError{Msg: msg, Err: err, Code: http.StatusNotFound})
		return
	}

	switch r.Method {
	case "GET":
		res, aerr := rc.Read(id)
		if aerr != nil {
			HTTPError(w, r, aerr)
			return
		}

		aerr = HTTPSendJSON(w, res, http.StatusOK)
		if aerr != nil {
			log.Println("Couldn't convert your error to JSON.")
			return
		}

	case "PUT":
		// Update Contender from json in r.Body and rc.Update()
		msg := "UPDATE is not exposed yet."
		HTTPError(w, r, &ApplicationError{Msg: msg, Code: http.StatusNotImplemented})

	case "DELETE":
		// rc.Destroy() by id
		msg := "DELETE is not exposed yet."
		HTTPError(w, r, &ApplicationError{Msg: msg, Code: http.StatusNotImplemented})

	default:
		// Unsupported Method, send Allow header back in error response
		w.Header().Add("Allow", "GET, PUT, DELETE")
		msg := fmt.Sprintf("Method \"%s\" is not supported at %s.", r.Method, r.URL.Path)
		HTTPError(w, r, &ApplicationError{Msg: msg, Code: http.StatusMethodNotAllowed})
		return
	}
}

// ServeResourceCollection redirects requests involving the collection.
func ServeResourceCollection(w http.ResponseWriter, r *http.Request, rc ResourceController, m models.Resource) {
	switch r.Method {
	case "GET":
		res, aerr := rc.ReadCollection()
		if aerr != nil {
			HTTPError(w, r, aerr)
			return
		}

		aerr = HTTPSendJSON(w, res, http.StatusOK)
		if aerr != nil {
			log.Println("Couldn't convert your error to JSON.")
			return
		}

	case "POST":
		// Create Contender from json in r.Body and rc.Create()
		msg := "CREATE is not exposed yet."
		HTTPError(w, r, &ApplicationError{Msg: msg, Code: http.StatusNotImplemented})

	case "DELETE":
		// Query for all contenders then rc.Destroy()
		msg := "DELETE is not exposed yet."
		HTTPError(w, r, &ApplicationError{Msg: msg, Code: http.StatusNotImplemented})

	default:
		// Unsupported Method, send Allow header back in error response
		w.Header().Add("Allow", "GET, POST, DELETE")
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
	log.Printf("Error: %s\n", aerr.Msg)

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
