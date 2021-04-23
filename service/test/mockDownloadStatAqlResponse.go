package test

import (
	"github.com/chanti529/repostats/util"
	"time"
)

func GetDownloadStatMockAqlResponse() util.AqlResult {

	now := time.Now()

	return util.AqlResult{
		Results: []*util.AqlItem{
			{
				Repo:       "repo1",
				Path:       "folder1",
				Name:       "file1",
				Modified:   now,
				ModifiedBy: "user1",
				Stats: []util.AqlItemStats{
					{
						Downloads:  10,
						Downloaded: now,
					},
				},
			},
			{
				Repo:       "repo1",
				Path:       "folder1",
				Name:       "file2",
				Modified:   now,
				ModifiedBy: "user2",
				Stats: []util.AqlItemStats{
					{
						Downloads:  9,
						Downloaded: now,
					},
				},
			},
			{
				Repo:       "repo1",
				Path:       "folder2",
				Name:       "file1",
				Modified:   now,
				ModifiedBy: "user1",
				Stats: []util.AqlItemStats{
					{
						Downloads:  8,
						Downloaded: now,
					},
				},
			},
			{
				Repo:       "repo1",
				Path:       "folder2",
				Name:       "file2",
				Modified:   now,
				ModifiedBy: "user2",
				Stats: []util.AqlItemStats{
					{
						Downloads:  7,
						Downloaded: now,
					},
				},
			},
			{
				Repo:       "repo2",
				Path:       "folder1",
				Name:       "file1",
				Modified:   now,
				ModifiedBy: "user1",
				Stats: []util.AqlItemStats{
					{
						Downloads:  6,
						Downloaded: now,
					},
				},
			},
			{
				Repo:       "repo2",
				Path:       "folder1",
				Name:       "file2",
				Modified:   now,
				ModifiedBy: "user2",
				Stats: []util.AqlItemStats{
					{
						Downloads:  5,
						Downloaded: now,
					},
				},
			},
			{
				Repo:       "repo2",
				Path:       "folder2",
				Name:       "file1",
				Modified:   now,
				ModifiedBy: "user1",
				Stats: []util.AqlItemStats{
					{
						Downloads:  4,
						Downloaded: now,
					},
				},
			},
			{
				Repo:       "repo2",
				Path:       "folder2",
				Name:       "file2",
				Modified:   now,
				ModifiedBy: "user2",
				Stats: []util.AqlItemStats{
					{
						Downloads:  3,
						Downloaded: now,
					},
				},
			},
			{
				Repo:       "anotherrepo",
				Path:       "folder",
				Name:       "file",
				Modified:   now,
				ModifiedBy: "user3",
				Stats: []util.AqlItemStats{
					{
						Downloads:  5,
						Downloaded: now,
					},
				},
			},
		},
	}
}
