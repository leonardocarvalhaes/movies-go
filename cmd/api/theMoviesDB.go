package main

import (
	"backend/internal/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type TheMovieDBResponse struct {
	Page         int `json:"page"`
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_results"`
	Results      []struct {
		ID          int     `json:"id"`
		Title       string  `json:"title"`
		ReleaseDate string  `json:"release_date"`
		VoteCount   int     `json:"vote_count"`
		VoteAverage float32 `json:"vote_average"`
		Overview    string  `json:"overview"`
		PosterPath  string  `json:"poster_path"`
	} `json:"results"`
}

func (app *application) downloadMovieData(movie *models.Movie) *models.Movie {
	tmdbResponse, err := app.searchTMDB(movie.Title)

	if err != nil {
		return movie
	}

	if len(tmdbResponse.Results) > 0 {
		movie.Image = tmdbResponse.Results[0].PosterPath
		movie.Description = tmdbResponse.Results[0].Overview
	}

	return movie
}

func (app *application) searchTMDB(searchTerm string) (*TheMovieDBResponse, error) {
	client := &http.Client{}
	uri := fmt.Sprintf("https://api.themoviedb.org/3/search/movie?api_key=%s", app.TMDBAPIKey)

	request, err := http.NewRequest("GET", uri+"&query="+url.QueryEscape(searchTerm), nil)

	if err != nil {
		return nil, err
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")

	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	var responseObject TheMovieDBResponse

	json.Unmarshal(bodyBytes, &responseObject)

	return &responseObject, nil
}
