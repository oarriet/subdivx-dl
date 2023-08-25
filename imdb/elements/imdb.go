package elements

type ImdbMovie struct {
	ID               string            `json:"id,omitempty"`
	ReviewAPIPath    string            `json:"review_api_path,omitempty"`
	Imdb             string            `json:"imdb,omitempty"`
	ContentType      string            `json:"contentType,omitempty"`
	ProductionStatus string            `json:"productionStatus,omitempty"`
	Title            string            `json:"title,omitempty"`
	Image            string            `json:"image,omitempty"`
	Images           []string          `json:"images,omitempty"`
	Plot             string            `json:"plot,omitempty"`
	Rating           Rating            `json:"rating,omitempty"`
	Award            Award             `json:"award,omitempty"`
	ContentRating    string            `json:"contentRating,omitempty"`
	Genre            []string          `json:"genre,omitempty"`
	ReleaseDetailed  ReleaseDetailed   `json:"releaseDetailed,omitempty"`
	Year             int               `json:"year,omitempty"`
	SpokenLanguages  []SpokenLanguages `json:"spokenLanguages,omitempty"`
	FilmingLocations []interface{}     `json:"filmingLocations,omitempty"`
	Runtime          string            `json:"runtime,omitempty"`
	RuntimeSeconds   int               `json:"runtimeSeconds,omitempty"`
	Actors           []string          `json:"actors,omitempty"`
	Directors        []interface{}     `json:"directors,omitempty"`
	TopCredits       []TopCredits      `json:"top_credits,omitempty"`
	Seasons          []Seasons         `json:"seasons,omitempty"`
	AllSeasons       []AllSeasons      `json:"all_seasons,omitempty"`
}

type Rating struct {
	Count int     `json:"count,omitempty"`
	Star  float64 `json:"star,omitempty"`
}

type Award struct {
	Wins        int `json:"wins,omitempty"`
	Nominations int `json:"nominations,omitempty"`
}

type ReleaseLocation struct {
	Country string `json:"country,omitempty"`
	Cca2    string `json:"cca2,omitempty"`
}

type OriginLocations struct {
	Country string `json:"country,omitempty"`
	Cca2    string `json:"cca2,omitempty"`
}

type ReleaseDetailed struct {
	Day             int               `json:"day,omitempty"`
	Month           int               `json:"month,omitempty"`
	Year            int               `json:"year,omitempty"`
	ReleaseLocation ReleaseLocation   `json:"releaseLocation,omitempty"`
	OriginLocations []OriginLocations `json:"originLocations,omitempty"`
}

type SpokenLanguages struct {
	Language string `json:"language,omitempty"`
	ID       string `json:"id,omitempty"`
}

type TopCredits struct {
	Name  string   `json:"name,omitempty"`
	Value []string `json:"value,omitempty"`
}

type Episodes struct {
	Idx           int    `json:"idx,omitempty"`
	No            string `json:"no,omitempty"`
	Title         string `json:"title,omitempty"`
	Image         string `json:"image,omitempty"`
	ImageLarge    string `json:"image_large,omitempty"`
	Plot          string `json:"plot,omitempty"`
	PublishedDate string `json:"publishedDate,omitempty"`
	Rating        Rating `json:"rating,omitempty"`
}

type Seasons struct {
	ID       string     `json:"id,omitempty"`
	APIPath  string     `json:"api_path,omitempty"`
	Name     string     `json:"name,omitempty"`
	Episodes []Episodes `json:"episodes,omitempty"`
}

type AllSeasons struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	APIPath string `json:"api_path,omitempty"`
}
