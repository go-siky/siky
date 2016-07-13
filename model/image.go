package model

import "time"

// Image of docker repository
type Image struct {
	Author        string    `json:"author"` //author
	Created       time.Time `json:"created"`
	DockerVersion string    `json:"docker_version"`
	OS            string    `json:"os"`
	ID            string    `json:"id"`
	Version       string    `json:"version"`
}

//ImageList struct implement sortable
type ImageList []Image

func (p ImageList) Len() int {
	return len(p)
}

func (p ImageList) Less(i, j int) bool {
	return p[i].Created.After(p[j].Created)
}

func (p ImageList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type Repository struct {
	Group  string   `json:"group"`
	Images []string `json:"images"`
}
