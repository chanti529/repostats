package util

import (
	"encoding/json"
	"fmt"
	"time"
)

type AqlSearchCriteria struct {
	Repos              []string
	PropertyFilter     []KeyValuePair
	ModifiedFrom       time.Time
	ModifiedTo         time.Time
	LastDownloadedFrom time.Time
	LastDownloadedTo   time.Time
}

func (c *AqlSearchCriteria) GetJson() ([]byte, error) {
	criteriaItems := make(map[string]interface{})

	var reposCriteriaOr []map[string]interface{}
	for _, item := range c.Repos {
		reposCriteriaOr = append(reposCriteriaOr, map[string]interface{}{
			"repo": map[string]string{
				"$match": item,
			},
		})
	}
	criteriaItems["$or"] = reposCriteriaOr

	for _, item := range c.PropertyFilter {
		criteriaItems[fmt.Sprintf("@%s", item.Key)] = map[string]string{
			"$match": item.Value,
		}
	}

	var timeBasedCriteria []map[string]interface{}
	if !c.ModifiedFrom.IsZero() {
		timeBasedCriteria = append(timeBasedCriteria, map[string]interface{}{
			"modified": map[string]time.Time{
				"$gt": c.ModifiedFrom,
			},
		})
	}

	if !c.ModifiedTo.IsZero() {
		timeBasedCriteria = append(timeBasedCriteria, map[string]interface{}{
			"modified": map[string]time.Time{
				"$lt": c.ModifiedTo,
			},
		})
	}

	if !c.LastDownloadedFrom.IsZero() {
		timeBasedCriteria = append(timeBasedCriteria, map[string]interface{}{
			"stat.downloaded": map[string]time.Time{
				"$gt": c.LastDownloadedFrom,
			},
		})
	}

	if !c.LastDownloadedTo.IsZero() {
		timeBasedCriteria = append(timeBasedCriteria, map[string]interface{}{
			"stat.downloaded": map[string]time.Time{
				"$lt": c.LastDownloadedTo,
			},
		})
	}

	if len(timeBasedCriteria) > 0 {
		criteriaItems["$and"] = timeBasedCriteria
	}

	return json.Marshal(criteriaItems)
}

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
