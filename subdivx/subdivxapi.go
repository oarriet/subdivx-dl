package subdivx

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/oarriet/subdivx-dl/subdivx/elements"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	subdivxAPIUrl     = "https://subdivx.com"
	subdivxAPITimeout = 30 * time.Second
)

type API interface {
	GetMoviesByTitle(title string) ([]elements.SubdivxMovie, error)
	DownloadSubtitle(pageUrl string) (io.ReadCloser, error)
	SaveSubtitle(subtitleReadCloser io.ReadCloser, filename string) error
}

type api struct {
}

func NewAPI() API {
	return &api{}
}

func (a *api) GetMoviesByTitle(title string) ([]elements.SubdivxMovie, error) {
	if len(title) == 0 {
		return nil, errors.New("title cannot be empty")
	}

	client := http.Client{
		Timeout: subdivxAPITimeout,
	}

	subdivxResponse, err := client.PostForm(subdivxAPIUrl,
		url.Values{
			"buscar2": {title},
			"accion":  {"5"},
			"masdesc": {""},
			"oxdown":  {"1"},
			"pg":      {"1"}, //TODO: pagination pageNum
		},
	)
	if err != nil {
		return nil, err
	}
	defer subdivxResponse.Body.Close()

	if subdivxResponse.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("subdivx API response status code: %d", subdivxResponse.StatusCode))
	}

	document, err := goquery.NewDocumentFromReader(subdivxResponse.Body)
	if err != nil {
		return nil, err
	}

	movies := make([]elements.SubdivxMovie, document.Find("div#menu_detalle_buscador").Length())

	document.Find("div#menu_detalle_buscador").Each(func(i int, selection *goquery.Selection) {
		if len(selection.Text()) == 0 {
			movies[i].Title = "No title"
		} else {
			movies[i].Title = selection.Text()
		}
		movies[i].Url, _ = selection.Find("a").Attr("href")
	})

	document.Find("div#buscador_detalle_sub").Each(func(i int, selection *goquery.Selection) {
		if len(selection.Text()) == 0 {
			movies[i].Description = "No description"
		} else {
			movies[i].Description = selection.Text()
		}
	})

	document.Find("div#buscador_detalle_sub_datos").Each(func(i int, selection *goquery.Selection) {
		downloadsCount, cds, commentsCount, format, uploadedBy, uploadedDate := stripData(selection.Text())

		movies[i].DownloadsCount = downloadsCount
		movies[i].Cds = cds
		movies[i].CommentsCount = commentsCount
		movies[i].Format = format
		movies[i].UploadedBy = uploadedBy
		movies[i].UploadedDate = uploadedDate
	})

	return movies, nil
}

/*
	stripData removes all the html tags from the data string

This is an example of the data string:

	"Downloads: 45,376 Cds: 1 Comentarios: 30 Formato: SubRip Subido por: FixXxer_mt  el 01/11/2007"
*/
func stripData(data string) (downloadsCount int, cds int, commentsCount int, format string, uploadedBy string, uploadedDate time.Time) {
	// Regular expression to match the desired components
	re := regexp.MustCompile(`Downloads:\s*([\d,]+)\s*Cds:\s*(\d+)\s*Comentarios:\s*(\d+)\s*Formato:\s*([\w\s]+)\s*Subido por:\s*([\w_]+)\s*el\s*(\d{2}/\d{2}/\d{4})`)

	// Find the matched parts in the input data
	matches := re.FindStringSubmatch(data)

	if len(matches) == 7 {
		downloadsStr := strings.ReplaceAll(matches[1], ",", "")
		downloadsCount, _ = strconv.Atoi(downloadsStr)
		cds, _ = strconv.Atoi(matches[2])
		commentsCount, _ = strconv.Atoi(matches[3])
		format = matches[4]
		uploadedBy = matches[5]

		uploadedDateStr := strings.ReplaceAll(matches[6], "/", "-")
		uploadedDate, _ = time.Parse("02-01-2006", uploadedDateStr)
	}

	return downloadsCount, cds, commentsCount, format, uploadedBy, uploadedDate
}

// DownloadSubtitle returns the subtitle file from the given downloadPageUrl, caller must close the io.ReadCloser
func (a *api) DownloadSubtitle(downloadPageUrl string) (io.ReadCloser, error) {
	if len(downloadPageUrl) == 0 {
		return nil, errors.New("downloadPageUrl cannot be empty")
	}

	client := http.Client{
		Timeout: subdivxAPITimeout,
	}

	downloadPageResponse, err := client.Get(downloadPageUrl)
	if err != nil {
		return nil, err
	}
	defer downloadPageResponse.Body.Close()

	if downloadPageResponse.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("downloadPageResponse status code: %d", downloadPageResponse.StatusCode))
	}

	document, err := goquery.NewDocumentFromReader(downloadPageResponse.Body)
	if err != nil {
		return nil, err
	}

	downloadLink, found := document.Find(".link1").Attr("href")
	if !found {
		return nil, errors.New("download link not found")
	}

	subdivxURL, err := url.Parse(subdivxAPIUrl)
	if err != nil {
		return nil, err
	}

	subdivxURL.Path = path.Join(subdivxURL.Path, downloadLink)

	downloadResponse, err := client.Get(subdivxURL.String())
	if err != nil {
		return nil, err
	}
	defer downloadResponse.Body.Close()

	if downloadResponse.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("downloadResponse status code: %d", downloadResponse.StatusCode))
	}

	redirectUrl := downloadResponse.Request.URL.Scheme + "://" + downloadResponse.Request.URL.Host + downloadResponse.Request.URL.Path

	redirectResponse, err := client.Get(redirectUrl)
	if err != nil {
		return nil, err
	}
	defer redirectResponse.Body.Close()

	if redirectResponse.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("redirectResponse status code: %d", redirectResponse.StatusCode))
	}

	return downloadResponse.Body, nil
}

// SaveSubtitle saves the subtitle file to the given filename. This func will close the subdivxSubtitle io.ReadCloser
func (a *api) SaveSubtitle(subtitleReadCloser io.ReadCloser, filename string) error {
	defer subtitleReadCloser.Close()

	if len(filename) == 0 {
		return errors.New("filename cannot be empty")
	}

	//for now let's create it under the current directory
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, subtitleReadCloser)
	if err != nil {
		return err
	}

	return nil
}
