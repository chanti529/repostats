package util

import (
	"fmt"
	"time"
)

type AqlResult struct {
	Results []*AqlItem `json:"results"`
}

type AqlItem struct {
	Repo       string         `json:"repo"`
	Path       string         `json:"path"`
	Name       string         `json:"name"`
	Size       int            `json:"size"`
	Modified   time.Time      `json:"modified"`
	ModifiedBy string         `json:"modified_by"`
	Stats      []aqlItemStats `json:"stats"`
}

func (i *AqlItem) GetFullPath() string {
	if i.Path != "." {
		return fmt.Sprintf("%s/%s/%s", i.Repo, i.Path, i.Name)
	} else {
		return fmt.Sprintf("%s/%s", i.Repo, i.Name)
	}
}

type aqlItemStats struct {
	Downloads  int       `json:"downloads"`
	Downloaded time.Time `json:"downloaded"`
}
