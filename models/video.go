package models

type Invitation struct {
	ProjectID  string       `json:"projectId"`
	Background Background   `json:"background"`
	Audio      *AudioConfig `json:"audio,omitempty"`
	Slides     []Slide      `json:"slides"`
}

type Background struct {
	URL string `json:"url"`
}

type AudioConfig struct {
	URL               string  `json:"url"`
	Volume            float64 `json:"volume,omitempty"`
	Loop              bool    `json:"loop,omitempty"`
	FadeIn            float64 `json:"fadeIn,omitempty"`
	FadeOut           float64 `json:"fadeOut,omitempty"`
	MixWithBackground bool    `json:"mixWithBackground,omitempty"`
	BackgroundVolume  float64 `json:"backgroundVolume,omitempty"`
}

type Slide struct {
	SlideID         int         `json:"slideId"`
	DisplayDuration int         `json:"displayDuration"`
	Start           int         `json:"-"`
	End             int         `json:"-"`
	Transition      *Transition `json:"transition,omitempty"`
	Elements        []Element   `json:"elements"`
}

type Transition struct {
	Type     string  `json:"type"`
	Duration float64 `json:"duration"`
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

	Animation *Animation `json:"animation,omitempty"`
}

type Animation struct {
	Type     string  `json:"type"`
	Duration float64 `json:"duration,omitempty"`
	Distance int     `json:"distance,omitempty"`
}
