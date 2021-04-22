package service

import (
	"bytes"
	"encoding/json"
	"github.com/chanti529/repostats/service/test"
	"github.com/jfrog/jfrog-client-go/artifactory"
	"io"
	"io/ioutil"
	"reflect"
	"regexp"
	"testing"
)

type ServiceManagerSizeMock struct {
	artifactory.EmptyArtifactoryServicesManager
}

func (smm *ServiceManagerSizeMock) Aql(aql string) (io.ReadCloser, error) {
	aqlResult := test.GetSizeStatMockAqlResponse()
	responseJson, err := json.Marshal(aqlResult)
	if err != nil {
		return nil, err
	}

	return ioutil.NopCloser(bytes.NewReader(responseJson)), nil
}

func TestGetSizeStat(t *testing.T) {

	testCases := []struct {
		conf           RepoStatConfiguration
		expectedResult []StatItem
	}{
		/*
			Validate top 1 most downloaded
		*/
		{
			conf: RepoStatConfiguration{
				PageSize:             1000,
				MaxConcurrentWorkers: 1,
				Type:                 "artifact",
				Sort:                 "desc",
				Limit:                1,
			},
			expectedResult: []StatItem{
				{
					Id:    "repo1/folder1/file1",
					Value: 10,
				},
			},
		},
		/*
			Validate top 1 least downloaded
		*/
		{
			conf: RepoStatConfiguration{
				PageSize:             1000,
				MaxConcurrentWorkers: 1,
				Type:                 "artifact",
				Sort:                 "asc",
				Limit:                1,
			},
			expectedResult: []StatItem{
				{
					Id:    "repo2/folder2/file2",
					Value: 3,
				},
			},
		},
		/*
			Validate ranking response (top 3)
		*/
		{
			conf: RepoStatConfiguration{
				PageSize:             100,
				MaxConcurrentWorkers: 1,
				Type:                 "artifact",
				Sort:                 "desc",
				Limit:                3,
			},
			expectedResult: []StatItem{
				{
					Id:    "repo1/folder1/file1",
					Value: 10,
				},
				{
					Id:    "repo1/folder1/file2",
					Value: 9,
				},
				{
					Id:    "repo1/folder2/file1",
					Value: 8,
				},
			},
		},
		/*
			Validate alpha numeric sorting
		*/
		{
			conf: RepoStatConfiguration{
				PageSize:             100,
				MaxConcurrentWorkers: 1,
				Type:                 "artifact",
				Sort:                 "alpha",
				Limit:                3,
			},
			expectedResult: []StatItem{
				{
					Id:    "anotherrepo/folder/file",
					Value: 5,
				},
				{
					Id:    "repo1/folder1/file1",
					Value: 10,
				},
				{
					Id:    "repo1/folder1/file2",
					Value: 9,
				},
			},
		},
		/*
			Validate most downloaded folder
		*/
		{
			conf: RepoStatConfiguration{
				PageSize:             1000,
				MaxConcurrentWorkers: 1,
				Type:                 "folder",
				MaxDepth:             2,
				Sort:                 "desc",
				Limit:                1,
			},
			expectedResult: []StatItem{
				{
					Id:    "repo1/folder1",
					Value: 19,
				},
			},
		},
		/*
			Validate most downloaded repo
		*/
		{
			conf: RepoStatConfiguration{
				PageSize:             1000,
				MaxConcurrentWorkers: 1,
				Type:                 "repo",
				Sort:                 "desc",
				Limit:                1,
			},
			expectedResult: []StatItem{
				{
					Id:    "repo1",
					Value: 34,
				},
			},
		},
		/*
			Validate most downloaded per user
		*/
		{
			conf: RepoStatConfiguration{
				PageSize:             1000,
				MaxConcurrentWorkers: 1,
				Type:                 "user",
				Sort:                 "desc",
				Limit:                1,
			},
			expectedResult: []StatItem{
				{
					Id:    "user1",
					Value: 28,
				},
			},
		},
		/*
			Validate path filter
		*/
		{
			conf: RepoStatConfiguration{
				PageSize:             1000,
				MaxConcurrentWorkers: 1,
				Type:                 "artifact",
				Sort:                 "desc",
				Limit:                1,
				FilterPathRegexp:     regexp.MustCompile("repo2.*"),
			},
			expectedResult: []StatItem{
				{
					Id:    "repo2/folder1/file1",
					Value: 6,
				},
			},
		},
		/*
			Validate multi thread processing
		*/
		{
			conf: RepoStatConfiguration{
				PageSize:             1,
				MaxConcurrentWorkers: 5,
				Type:                 "artifact",
				Sort:                 "desc",
				Limit:                1,
			},
			expectedResult: []StatItem{
				{
					Id:    "repo1/folder1/file1",
					Value: 10,
				},
			},
		},
	}

	serviceManagerMock := ServiceManagerSizeMock{}

	for _, testCase := range testCases {
		result, err := GetSizeStat(&testCase.conf, &serviceManagerMock)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(result, testCase.expectedResult) {
			t.Errorf("Got wrong size stat result for conf %+v.\nGot: %+v\n"+
				"Exp: %+v", testCase.conf, result, testCase.expectedResult)
		}
	}
}
