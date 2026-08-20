package models

type Invitation struct {
	ProjectID  string     `json:"projectId"`
	Background Background `json:"background"`
	Slides     []Slide    `json:"slides"`
}

type Background struct {
	URL string `json:"url"`
}

type Slide struct {
	SlideID         int       `json:"slideId"`
	DisplayDuration int       `json:"displayDuration"`
	Start           int       `json:"-"`
	End             int       `json:"-"`
	Elements        []Element `json:"elements"`
}

type Element struct {
	ID   string `json:"id"`
	Type string `json:"type"`

	Text      string `json:"text,omitempty"`
	AssetPath string `json:"assetPath,omitempty"`

	X int `json:"x"`
	Y int `json:"y"`

	Width  int `json:"width,omitempty"`
	Height int `json:"height,omitempty"`

	FontSize  int    `json:"fontSize,omitempty"`
	FontColor string `json:"fontColor,omitempty"`
	FontFile  string `json:"fontFile,omitempty"`

	Align string `json:"align,omitempty"`

	Opacity  float64 `json:"opacity,omitempty"`
	Rotation float64 `json:"rotation,omitempty"`
}
