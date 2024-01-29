package main

import (
	"backend/internal/dtos"
	"backend/internal/models"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func (app *application) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)

	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)

	if err != nil {
		return err
	}

	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1024 * 1024

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)

	dec.DisallowUnknownFields()

	err := dec.Decode(data)

	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})

	if err != io.EOF {
		return errors.New("body must contain a single JSON")
	}

	return nil
}

func (app *application) errorJSON(w http.ResponseWriter, err error, status ...int) {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload dtos.JSONResponse
	payload.Error = true
	payload.Message = err.Error()

	app.writeJSON(w, statusCode, payload)
}

func (app *application) fromRequestToMovie(w http.ResponseWriter, r *http.Request) (*models.Movie, error) {
	var movie models.Movie

	err := app.readJSON(w, r, &movie)

	if err != nil {
		return nil, err
	}

	app.downloadMovieData(&movie)

	return &movie, nil
}
