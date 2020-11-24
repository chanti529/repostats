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
	// We cannot paginate this query since it relies on fields from stat subdomain
	aqlDownloadTemplate = `items.find({
			"repo": "%s" 
		}).include("repo", "path", "name", "modified", "modified_by", "stat.downloads", "stat.downloaded")`
)

func GetDownloadStat(conf *RepoStatConfiguration) ([]StatItem, error) {
	servicesManager, err := utils.CreateServiceManager(conf.RtDetails, false)
	if err != nil {
		return nil, err
	}

	aql := fmt.Sprintf(aqlDownloadTemplate, conf.Repos[0])

	reader, err := servicesManager.Aql(aql)
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

	itemsCount := len(parsedResult.Results)

	// TODO: Make page size configurable
	pageSize := 50000

	// TODO: Make number of workers configurable
	numberOfWorkers := 5
	workersLock := make(chan bool, numberOfWorkers)

	var mapperWorkers []*statMapper
	getValueFunc := func(item *util.AqlItem) int {
		return item.Stats[0].Downloads
	}

	scheduledItems := 0
	/*
		While there are new items, filter and map them to the requested identity
	*/
	for scheduledItems < itemsCount {

		mapper := newStatMapper()
		mapper.GetValueFunc = getValueFunc

		mapperWorkers = append(mapperWorkers, mapper)
		pageInitialIndex := scheduledItems
		pageFinalIndex := scheduledItems + pageSize
		if pageFinalIndex > itemsCount {
			pageFinalIndex = itemsCount
		}

		// Start mappers in parallel
		workersLock <- true
		go func(m *statMapper, initialIndex, finalIndex int) {
			m.process(parsedResult.Results[initialIndex:finalIndex], conf)
			<-workersLock
		}(mapper, pageInitialIndex, pageFinalIndex)

		scheduledItems = scheduledItems + pageSize
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
