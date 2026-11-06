package mal

// Picture represents image URLs for an anime or manga.
type Picture struct {
	Medium string `json:"medium"`
	Large  string `json:"large,omitempty"`
}

// AlternativeTitles represents localized and alternative titles.
type AlternativeTitles struct {
	Synonyms []string `json:"synonyms,omitempty"`
	En       string   `json:"en,omitempty"`
	Ja       string   `json:"ja,omitempty"`
}

// Genre represents a media genre.
type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// SeasonInfo represents a broadcast season.
type SeasonInfo struct {
	Year   int    `json:"year"`
	Season string `json:"season"`
}

// Broadcast represents the scheduled broadcast time in JST.
type Broadcast struct {
	DayOfTheWeek string `json:"day_of_the_week,omitempty"`
	StartTime    string `json:"start_time,omitempty"`
}

// Studio represents an animation studio.
type Studio struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Paging represents pagination links returned in list endpoints.
type Paging struct {
	Previous string `json:"previous,omitempty"`
	Next     string `json:"next,omitempty"`
}
