package service

import (
	"encoding/json"
	"fmt"
	"github.com/chanti529/jfrog-cli-plugin-template/util"
	"github.com/jfrog/jfrog-cli-core/artifactory/utils"
	"io/ioutil"
)

const (
	//TODO: Add additional filters to AQL
	//TODO: Support multiple repos
	aqlSizeTemplate = `items.find({
			"repo": "%s" 
		}).include("repo", "path", "name", "created", "modified", "modified_by", "size").sort({
			"$asc":["created"]
		})`
)

func GetSizeStat(conf *RepoStatConfiguration) ([]StatItem, error) {
	servicesManager, err := utils.CreateServiceManager(conf.RtDetails, false)
	if err != nil {
		return nil, err
	}

	aql := fmt.Sprintf(aqlSizeTemplate, conf.Repos[0])

	// TODO: Make page size configurable
	pageSize := 500

	// TODO: Make number of workers configurable
	numberOfWorkers := 5
	workersLock := make(chan bool, numberOfWorkers)

	var mapperWorkers []*statMapper
	getValueFunc := func(item *util.AqlItem) int {
		return item.Size
	}

	itemsCount := 0
	// Set itemsInPage with pageSize to force first page to be fetched
	itemsInPage := pageSize
	for itemsInPage == pageSize {

		// Query page results
		pageAql := fmt.Sprintf("%s.offset(%v).limit(%v)", aql, itemsCount, pageSize)
		reader, err := servicesManager.Aql(pageAql)
		if err != nil {
			return nil, err
		}
		defer reader.Close()
		result, err := ioutil.ReadAll(reader)
		if err != nil {
			return nil, err
		}

		parsedResult := new(util.AqlResult)
		if err = json.Unmarshal(result, parsedResult); err != nil {
			return nil, err
		}

		itemsInPage = len(parsedResult.Results)
		itemsCount = itemsCount + itemsInPage

		if itemsInPage > 0 {
			mapper := newStatMapper()
			mapper.GetValueFunc = getValueFunc
			mapperWorkers = append(mapperWorkers, mapper)

			// Start mappers in parallel
			workersLock <- true
			go func(m *statMapper, items []*util.AqlItem) {
				m.process(items, conf)
				<-workersLock
			}(mapper, parsedResult.Results)
		}
	}

	// Wait for mappers to finish
	err = waitForWorkers(mapperWorkers)
	if err != nil {
		return nil, err
	}

	/*
		Reduce results from workers to their identity and apply limits
	*/
	repoStatsResult, err := reduce(mapperWorkers, conf)
	if err != nil {
		return nil, err
	}
	return repoStatsResult, nil
}
