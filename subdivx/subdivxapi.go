package subdivx

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gen2brain/go-unarr"
	"github.com/google/uuid"
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
	subdivxCookie     = "__cfduid=dea8419e3bf838c5ec1b8624c00ba126e1599785667; con_impr=5; cant_down=16; bajo_una_vez=0; bajo_una_vez_diario=0; contd=3; cs15=566391; cs14=215575; cs13=277494; __cf_bm=edba632a5a68f4ad6890f4f58ff571044dc84d1a-1599793485-1800-Ac0IxmDEnWMbTjrtkhEluRQMTH6hnt2KhSJGCa7KPLxY"
)

type API interface {
	GetMoviesByTitle(title string) ([]elements.SubdivxMovie, error)
	DownloadSubtitle(pageUrl string) (rc io.ReadCloser, contentType string, err error)
	SaveSubtitle(subtitleReadCloser io.ReadCloser, contentType string, folderPath string) ([]string, error)
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
func (a *api) DownloadSubtitle(downloadPageUrl string) (io.ReadCloser, string, error) {
	if len(downloadPageUrl) == 0 {
		return nil, "", errors.New("downloadPageUrl cannot be empty")
	}

	client := http.Client{
		Timeout: subdivxAPITimeout,
	}

	//add user agent
	req, err := http.NewRequest(http.MethodGet, downloadPageUrl, nil)
	if err != nil {
		return nil, "", err
	}

	downloadPageResponse, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer downloadPageResponse.Body.Close()

	if downloadPageResponse.StatusCode != http.StatusOK {
		return nil, "", errors.New(fmt.Sprintf("downloadPageResponse status code: %d", downloadPageResponse.StatusCode))
	}

	document, err := goquery.NewDocumentFromReader(downloadPageResponse.Body)
	if err != nil {
		return nil, "", err
	}

	downloadLink, found := document.Find(".link1").Attr("href")
	if !found {
		return nil, "", errors.New("download link not found")
	}

	subdivxURL, err := url.Parse(subdivxAPIUrl)
	if err != nil {
		return nil, "", err
	}

	subdivxURL.Path = path.Join(subdivxURL.Path, downloadLink)

	//add user agent
	req, err = http.NewRequest(http.MethodGet, subdivxURL.String(), nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Cookie", subdivxCookie)

	downloadResponse, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer downloadResponse.Body.Close()

	redirectUrl := downloadResponse.Request.URL.Scheme + "://" + downloadResponse.Request.URL.Host + downloadResponse.Request.URL.Path

	req, err = http.NewRequest(http.MethodGet, redirectUrl, nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Cookie", subdivxCookie)

	redirectResponse, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}

	return redirectResponse.Body, redirectResponse.Header.Get("Content-Type"), nil
}

// SaveSubtitle saves the subtitle file to the given filename. This func will close the subdivxSubtitle io.ReadCloser
func (a *api) SaveSubtitle(subtitleReadCloser io.ReadCloser, contentType string, folderPath string) ([]string, error) {
	defer subtitleReadCloser.Close()

	if len(folderPath) == 0 {
		return nil, errors.New("folderPath cannot be empty")
	}

	//create folder if it doesn't exist
	err := os.MkdirAll(path.Dir(folderPath), os.ModePerm)
	if err != nil {
		return nil, err
	}

	var extension string
	//var format archiver.Extractor
	if strings.EqualFold(contentType, "application/rar") {
		extension = ".rar"
	}
	if strings.EqualFold(contentType, "application/zip") {
		extension = ".zip"
	}

	uuid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	fileName := fmt.Sprintf("%s%s", uuid.String(), extension)

	//for now let's create it under the current directory
	file, err := os.Create(path.Join(folderPath, fileName))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	_, err = io.Copy(file, subtitleReadCloser)
	if err != nil {
		return nil, err
	}

	archive, err := unarr.NewArchive(file.Name())
	if err != nil {
		return nil, err
	}
	defer archive.Close()

	subNames, err := archive.Extract(folderPath)
	if err != nil {
		return nil, err
	} else {
		//delete the archive file
		_ = os.Remove(file.Name())
		//ignore the error because the whole point of this func is to save the subtitle
	}

	return subNames, nil
}
