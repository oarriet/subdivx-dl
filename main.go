package main

import (
	"github.com/oarriet/subdivx-dl/subdivx"
	"log"
)

const (
	folderToDownload = "build"

	movieToSearch = "tt0245574"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	////we get the movie name from the id
	//imdbAPI := imdb.NewAPI()
	//movie, err := imdbAPI.GetMovieById(movieToSearch)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//we get the subdivx data from the movie name
	subdivxAPI := subdivx.NewAPI()
	//subdivxMovies, err := subdivxAPI.GetMoviesByTitle(fmt.Sprintf("%s %d", movie.Title, movie.Year))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	////let's print the subdivx data
	//for _, subdivxMovie := range subdivxMovies {
	//	jMovie, _ := json.MarshalIndent(subdivxMovie, "", "\t")
	//	log.Println(string(jMovie))
	//}
	//
	////let's download the first subtitle
	//subdivxMovie := subdivxMovies[0]
	subdivxSubtitle, contentType, err := subdivxAPI.DownloadSubtitle("https://www.subdivx.com/X666XMzQ3MTk2X-y-tu-mam%EF%BF%BD-tambi%C3%A9n-2001.html")
	if err != nil {
		log.Fatal(err)
	}
	//defer subdivxSubtitle.Close()

	//save the subtitle
	err = subdivxAPI.SaveSubtitle(subdivxSubtitle, contentType, folderToDownload)
	if err != nil {
		log.Fatal(err)
	}
}
