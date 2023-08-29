package main

import (
	"encoding/json"
	"fmt"
	"github.com/oarriet/subdivx-dl/imdb"
	"github.com/oarriet/subdivx-dl/subdivx"
	"log"
	"path"
)

const (
	folderToDownload = "build"

	movieToSearch = "tt1001520"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//we get the movie name from the id
	imdbAPI := imdb.NewAPI()
	movie, err := imdbAPI.GetMovieById(movieToSearch)
	if err != nil {
		log.Fatal(err)
	}

	//we get the subdivx data from the movie name
	subdivxAPI := subdivx.NewAPI()
	subdivxMovies, err := subdivxAPI.GetMoviesByTitle(fmt.Sprintf("%s %d", movie.Title, movie.Year))
	if err != nil {
		log.Fatal(err)
	}

	//let's print the subdivx data
	for _, subdivxMovie := range subdivxMovies {
		jMovie, _ := json.MarshalIndent(subdivxMovie, "", "\t")
		log.Println(string(jMovie))
	}

	//let's download the first subtitle
	subdivxMovie := subdivxMovies[0]
	subdivxSubtitle, err := subdivxAPI.DownloadSubtitle(subdivxMovie.Url)
	if err != nil {
		log.Fatal(err)
	}
	//defer subdivxSubtitle.Close()

	//save the subtitle
	err = subdivxAPI.SaveSubtitle(subdivxSubtitle, path.Join(folderToDownload, fmt.Sprintf("%s.zip", subdivxMovie.Title)))
	if err != nil {
		log.Fatal(err)
	}
}
