package util

import (
	"encoding/json"
	"testing"
	"time"
)

func TestAqlItem_GetFullPath(t *testing.T) {
	testCases := []struct {
		aqlItem          AqlItem
		expectedFullPath string
	}{
		{
			aqlItem: AqlItem{
				Repo: "repo",
				Path: ".",
				Name: "name",
			},
			expectedFullPath: "repo/name",
		},
		{
			aqlItem: AqlItem{
				Repo: "repo",
				Path: "path",
				Name: "name",
			},
			expectedFullPath: "repo/path/name",
		},
	}

	for _, testCase := range testCases {
		fullPath := testCase.aqlItem.GetFullPath()
		if fullPath != testCase.expectedFullPath {
			t.Errorf("Got wrong full path for AqlItem %+v. Got: %s, "+
				"Expected: %s", testCase.aqlItem, fullPath,
				testCase.expectedFullPath)
		}
	}
}

func TestAqlSearchCriteria_GetJson(t *testing.T) {

	now := time.Now()

	testCases := []struct {
		aqlSearchCriteria AqlSearchCriteria
		expectedJson      string
	}{
		{
			aqlSearchCriteria: AqlSearchCriteria{
				Repos:              []string{"repo1", "repo2"},
				PropertyFilter:     nil,
				ModifiedFrom:       time.Time{},
				ModifiedTo:         time.Time{},
				LastDownloadedFrom: time.Time{},
				LastDownloadedTo:   time.Time{},
			},
			expectedJson: `{"$or":[{"repo":{"$match":"repo1"}},{"repo":{"$match":"repo2"}}]}`,
		},
		{
			aqlSearchCriteria: AqlSearchCriteria{
				Repos: []string{"repo1", "repo2"},
				PropertyFilter: []KeyValuePair{
					{
						Key:   "prop1",
						Value: "value1",
					},
					{
						Key:   "prop2",
						Value: "value2",
					},
				},
				ModifiedFrom:       time.Time{},
				ModifiedTo:         time.Time{},
				LastDownloadedFrom: time.Time{},
				LastDownloadedTo:   time.Time{},
			},
			expectedJson: `{"$or":[{"repo":{"$match":"repo1"}},{"repo":{"$match"` +
				`:"repo2"}}],"@prop1":{"$match":"value1"},"@prop2":{"$match":"value2"}}`,
		},
		{
			aqlSearchCriteria: AqlSearchCriteria{
				Repos: []string{"repo1", "repo2"},
				PropertyFilter: []KeyValuePair{
					{
						Key:   "prop1",
						Value: "value1",
					},
					{
						Key:   "prop2",
						Value: "value2",
					},
				},
				ModifiedFrom:       now,
				ModifiedTo:         time.Time{},
				LastDownloadedFrom: time.Time{},
				LastDownloadedTo:   time.Time{},
			},
			expectedJson: `{"$and":[{"modified":{"$gt":` + getJsonTimestamp(now) +
				`}}],"$or":[{"repo":{"$match":"repo1"}},{"repo":{"$match":"repo2"}}],` +
				`"@prop1":{"$match":"value1"},"@prop2":{"$match":"value2"}}`,
		},
		{
			aqlSearchCriteria: AqlSearchCriteria{
				Repos:              []string{"repo1", "repo2"},
				PropertyFilter:     nil,
				ModifiedFrom:       now,
				ModifiedTo:         time.Time{},
				LastDownloadedFrom: time.Time{},
				LastDownloadedTo:   time.Time{},
			},
			expectedJson: `{"$and":[{"modified":{"$gt":` + getJsonTimestamp(now) +
				`}}],"$or":[{"repo":{"$match":"repo1"}},{"repo":{"$match":"repo2"}}]}`,
		},
		{
			aqlSearchCriteria: AqlSearchCriteria{
				Repos:              []string{"repo1", "repo2"},
				PropertyFilter:     nil,
				ModifiedFrom:       now,
				ModifiedTo:         now,
				LastDownloadedFrom: time.Time{},
				LastDownloadedTo:   time.Time{},
			},
			expectedJson: `{"$and":[{"modified":{"$gt":` + getJsonTimestamp(now) +
				`}},{"modified":{"$lt":` + getJsonTimestamp(now) +
				`}}],"$or":[{"repo":{"$match":"repo1"}},{"repo":{"$match":"repo2"}}]}`,
		},
		{
			aqlSearchCriteria: AqlSearchCriteria{
				Repos:              []string{"repo1", "repo2"},
				PropertyFilter:     nil,
				ModifiedFrom:       time.Time{},
				ModifiedTo:         time.Time{},
				LastDownloadedFrom: now,
				LastDownloadedTo:   time.Time{},
			},
			expectedJson: `{"$and":[{"stat.downloaded":{"$gt":` + getJsonTimestamp(now) +
				`}}],"$or":[{"repo":{"$match":"repo1"}},{"repo":{"$match":"repo2"}}]}`,
		},
		{
			aqlSearchCriteria: AqlSearchCriteria{
				Repos:              []string{"repo1", "repo2"},
				PropertyFilter:     nil,
				ModifiedFrom:       time.Time{},
				ModifiedTo:         time.Time{},
				LastDownloadedFrom: time.Time{},
				LastDownloadedTo:   now,
			},
			expectedJson: `{"$and":[{"stat.downloaded":{"$lt":` + getJsonTimestamp(now) +
				`}}],"$or":[{"repo":{"$match":"repo1"}},{"repo":{"$match":"repo2"}}]}`,
		},
	}

	for _, testCase := range testCases {
		json, err := testCase.aqlSearchCriteria.GetJson()
		if err != nil {
			t.Error(err)
		}

		jsonString := string(json)

		if jsonString != testCase.expectedJson {
			t.Errorf("Got wrong json for AqlSearchCriteria %+v.\nGot: %s\n"+
				"Exp: %s", testCase.aqlSearchCriteria, json,
				testCase.expectedJson)
		}
	}
}

func getJsonTimestamp(time2 time.Time) string {
	jsonTime, _ := json.Marshal(time2)
	return string(jsonTime)
}
