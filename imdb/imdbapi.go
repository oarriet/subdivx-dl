package imdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/oarriet/subdivx-dl/imdb/elements"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"
)

const (
	imdbAPIUrl        = "https://imdb-api.projects.thetuhin.com"
	imdbAPITimeout    = 30 * time.Second
	imdbAPIParamTitle = "title"
)

type API interface {
	GetMovieById(id string) (*elements.ImdbMovie, error)
}

type api struct {
}

func NewAPI() API {
	return &api{}
}

// GetMovieById returns a movie by its id using the imdb-api, provided by https://imdb-api.projects.thetuhin.com
func (a *api) GetMovieById(id string) (*elements.ImdbMovie, error) {
	if len(id) == 0 {
		return nil, errors.New("id cannot be empty")
	}

	imdbAPIURL, err := url.Parse(imdbAPIUrl)
	if err != nil {
		return nil, err
	}

	imdbAPIURL.Path = path.Join(imdbAPIURL.Path, imdbAPIParamTitle, id)

	client := http.Client{
		Timeout: imdbAPITimeout,
	}
	imdbResponse, err := client.Get(imdbAPIURL.String())
	if err != nil {
		return nil, err
	}
	defer imdbResponse.Body.Close()

	if imdbResponse.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("imdb API response status code: %d", imdbResponse.StatusCode))
	}

	body, err := io.ReadAll(imdbResponse.Body)
	if err != nil {
		return nil, err
	}

	movie := elements.ImdbMovie{}
	err = json.Unmarshal(body, &movie)
	if err != nil {
		return nil, err
	}

	return &movie, nil
}
