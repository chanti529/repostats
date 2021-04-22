package test

import (
	"github.com/chanti529/repostats/util"
	"time"
)

func GetSizeStatMockAqlResponse() util.AqlResult {

	now := time.Now()

	return util.AqlResult{
		Results: []*util.AqlItem{
			{
				Repo:       "repo1",
				Path:       "folder1",
				Name:       "file1",
				Modified:   now,
				ModifiedBy: "user1",
				Size:       10,
			},
			{
				Repo:       "repo1",
				Path:       "folder1",
				Name:       "file2",
				Modified:   now,
				ModifiedBy: "user2",
				Size:       9,
			},
			{
				Repo:       "repo1",
				Path:       "folder2",
				Name:       "file1",
				Modified:   now,
				ModifiedBy: "user1",
				Size:       8,
			},
			{
				Repo:       "repo1",
				Path:       "folder2",
				Name:       "file2",
				Modified:   now,
				ModifiedBy: "user2",
				Size:       7,
			},
			{
				Repo:       "repo2",
				Path:       "folder1",
				Name:       "file1",
				Modified:   now,
				ModifiedBy: "user1",
				Size:       6,
			},
			{
				Repo:       "repo2",
				Path:       "folder1",
				Name:       "file2",
				Modified:   now,
				ModifiedBy: "user2",
				Size:       5,
			},
			{
				Repo:       "repo2",
				Path:       "folder2",
				Name:       "file1",
				Modified:   now,
				ModifiedBy: "user1",
				Size:       4,
			},
			{
				Repo:       "repo2",
				Path:       "folder2",
				Name:       "file2",
				Modified:   now,
				ModifiedBy: "user2",
				Size:       3,
			},
			{
				Repo:       "anotherrepo",
				Path:       "folder",
				Name:       "file",
				Modified:   now,
				ModifiedBy: "user3",
				Size:       5,
			},
		},
	}
}
