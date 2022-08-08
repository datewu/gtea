package router

import "net/http"

type interceptor4xx struct {
	// Interceptor  404/403 response
	origWriter                 http.ResponseWriter
	overridden                 bool
	notFound, methodNotAllowed http.HandlerFunc
}

func (i *interceptor4xx) WriteHeader(rc int) {
	switch rc {
	// case 500:
	// 	http.Error(i.origWriter, "Custom 500 message / content", 500)
	case http.StatusNotFound:
		i.notFound(i.origWriter, nil)
	case http.StatusMethodNotAllowed:
		i.methodNotAllowed(i.origWriter, nil)
	default:
		i.origWriter.WriteHeader(rc)
		return
	}
	// if the default case didn't execute (and return) we must have overridden the output
	i.overridden = true
}

func (i *interceptor4xx) Write(b []byte) (int, error) {
	if !i.overridden {
		return i.origWriter.Write(b)
	}
	// Return nothing if we've overriden the response.
	return 0, nil
}

func (i *interceptor4xx) Header() http.Header {
	return i.origWriter.Header()
}
