package main

import (
	"github.com/Ankumeah/JSBEE/smtp_relay/emails"

	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func group(
	parent *http.ServeMux,
	prefix string,
	middleware func(http.Handler) http.Handler,
) *http.ServeMux {
	child := http.NewServeMux()

	parent.Handle(prefix+"/", middleware(
		http.StripPrefix(prefix, child),
	))

	return child
}

func main() {
	if err := emails.InitTemplates(); err != nil {
		log.Fatalln(err.Error())
	}

	mux := http.NewServeMux()

	api := group(mux,
		"/api"+envVars["API_VERSION"]+"/",
		func(h http.Handler) http.Handler { return h },
	)

	api.HandleFunc("POST /email", createEmailJob)

	if err := http.ListenAndServe("0.0.0.0", mux); err != nil {
		log.Fatalf("Error while serving: %v\n", err)
	}
}

func createEmailJob(
	w http.ResponseWriter,
	r *http.Request,
) {
	var request struct {
		Type       string `json:"type"`
		User       string `json:"user"`
		Email      string `json:"email"`
		PaperTitle string `json:"paper_title"`
		Feedback   string `json:"feedback,omitempty"`
	}

	if err := json.NewDecoder(
		r.Body,
	).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(request.Type) == "" ||
		strings.TrimSpace(request.User) == "" ||
		strings.TrimSpace(request.Email) == "" ||
		strings.TrimSpace(request.PaperTitle) == "" {

		w.WriteHeader(http.StatusBadRequest)
		http.Error(w,
			"Got one or more non feedback empty field",
			http.StatusBadRequest,
		)
		return
	}

	// TODO: Add to queue
	panic("Unimplimented")
}
