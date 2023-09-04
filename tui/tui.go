package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/oarriet/subdivx-dl/imdb"
	"github.com/oarriet/subdivx-dl/subdivx"
	"github.com/oarriet/subdivx-dl/subdivx/elements"
	"github.com/oarriet/subdivx-dl/utils"
	"strings"
)

const (
	folderToDownload = "build"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func NewModel() tea.Model {
	return initialModel()
}

type (
	errMsg error
)

type model struct {
	spinner   spinner.Model
	textInput textinput.Model
	table     table.Model
	textarea  textarea.Model
	err       error

	inProgress bool

	subdivxMovies []elements.SubdivxMovie
}

func initialModel() model {
	spin := spinner.New()
	spin.Spinner = spinner.Dot
	spin.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("57"))

	ti := textinput.New()
	ti.Placeholder = "tt0111161"
	ti.Focus()
	ti.CharLimit = 20
	ti.Width = 20

	columns := []table.Column{
		{Title: "Title", Width: 50},
		{Title: "UploadedBy", Width: 20},
		{Title: "DownloadsCount", Width: 20},
		{Title: "UploadedDate", Width: 20},
	}

	var rows []table.Row

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(15),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	ta := textarea.New()
	ta.Placeholder = "Input the IMDb movie/TV id and press enter"
	ta.SetHeight(3)
	ta.SetWidth(100)

	return model{
		spinner:       spin,
		textInput:     ti,
		table:         t,
		textarea:      ta,
		err:           nil,
		subdivxMovies: make([]elements.SubdivxMovie, 0),
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if len(m.textInput.Value()) == 0 {
				return m, nil
			}
			m.inProgress = true
			return m, tea.Sequence(m.spinner.Tick, m.refreshTableWithData(m.textInput.Value()))
		case tea.KeyDown, tea.KeyUp:
			m.table, cmd = m.table.Update(msg)
			m.textarea.SetValue(fmt.Sprintf("%s", m.subdivxMovies[m.table.Cursor()].Description))
			return m, cmd
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyCtrlD:
			if len(m.table.Rows()) == 0 {
				return m, nil
			}
			m.inProgress = true
			return m, tea.Sequence(m.spinner.Tick, m.downloadSubtitle(m.subdivxMovies[m.table.Cursor()]))
		}

	case spinner.TickMsg:
		if !m.inProgress {
			return m, nil
		} else {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}

	case []elements.SubdivxMovie:
		var rows []table.Row
		for _, subdivxMovie := range msg {
			rows = append(rows, table.Row{
				subdivxMovie.Title,
				subdivxMovie.UploadedBy,
				fmt.Sprintf("%s", utils.FormatIntWithCommasAndPoints(subdivxMovie.DownloadsCount)),
				fmt.Sprintf("%s", subdivxMovie.UploadedDate.Format("2006-01-02")),
			})
		}
		if len(rows) > 0 {
			m.textarea.SetValue(fmt.Sprintf("%s", msg[0].Description))
		}
		m.subdivxMovies = msg
		//stop the spinner
		m.inProgress = false
		m.table.SetRows(rows)
		m.table.Focus()
		return m, nil

	// We handle errors just like any other message
	case errMsg:
		m.textarea.SetValue(fmt.Sprintf("Error: %s", msg))
		//stop the spinner
		m.inProgress = false
		m.table.SetRows([]table.Row{})
		m.err = msg
		return m, nil

	case SubMsg:
		m.textarea.SetValue(fmt.Sprintf("Subtitle downloaded successfully: %s", strings.Join(msg.SubNames, ", ")))
		//stop the spinner
		m.inProgress = false
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf(
		"IMDb Movie/TV id? \n\n%s %s\n%s\n%s\n\n%s",
		m.spinner.View(),
		m.textInput.View(),
		baseStyle.Render(m.table.View()),
		m.textarea.View(),
		"(d to download, esc to quit)",
	) + "\n"
}

func (m model) refreshTableWithData(imdbId string) tea.Cmd {
	return func() tea.Msg {
		imdbAPI := imdb.NewAPI()
		movie, err := imdbAPI.GetMovieById(imdbId)
		if err != nil {
			return errMsg(err)
		}

		//we get the subdivx data from the movie name
		subdivxAPI := subdivx.NewAPI()
		subdivxMovies, err := subdivxAPI.GetMoviesByTitle(fmt.Sprintf("%s %d", movie.Title, movie.Year))
		if err != nil {
			return errMsg(err)
		}

		return subdivxMovies
	}
}

func (m model) downloadSubtitle(movie elements.SubdivxMovie) tea.Cmd {
	return func() tea.Msg {
		subdivxAPI := subdivx.NewAPI()
		subdivxSubtitle, contentType, err := subdivxAPI.DownloadSubtitle(movie.Url)
		if err != nil {
			return errMsg(err)
		}

		//save the subtitle
		subName, err := subdivxAPI.SaveSubtitle(subdivxSubtitle, contentType, folderToDownload)
		if err != nil {
			return errMsg(err)
		}
		return SubMsg{
			Succeed:  true,
			SubNames: subName,
		}
	}
}
