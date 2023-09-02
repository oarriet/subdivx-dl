package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/oarriet/subdivx-dl/tui"
	"log"
	"os"
)

const (
	folderToDownload = "build"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	p := tea.NewProgram(tui.NewModel())

	_, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	////let's download the first subtitle
	//subdivxMovie := subdivxMovies[0]
	//subdivxSubtitle, contentType, err := subdivxAPI.DownloadSubtitle("https://www.subdivx.com/X666XMzQ3MTk2X-y-tu-mam%EF%BF%BD-tambi%C3%A9n-2001.html")
	//if err != nil {
	//	log.Fatal(err)
	//}
	////defer subdivxSubtitle.Close()
	//
	////save the subtitle
	//err = subdivxAPI.SaveSubtitle(subdivxSubtitle, contentType, folderToDownload)
	//if err != nil {
	//	log.Fatal(err)
	//}
}
