package service

import (
	"encoding/json"
	"fmt"
	"github.com/chanti529/jfrog-cli-plugin-template/util"
	"github.com/jfrog/jfrog-cli-core/artifactory/utils"
	"io/ioutil"
)

const (
	// We cannot paginate this query since it relies on fields from stat subdomain
	aqlDownloadTemplate = `items.find(%s).include("repo", "path", "name", "modified", "modified_by", "stat.downloads", "stat.downloaded")`
)

func GetDownloadStat(conf *RepoStatConfiguration) ([]StatItem, error) {
	servicesManager, err := utils.CreateServiceManager(conf.RtDetails, false)
	if err != nil {
		return nil, err
	}

	aqlCriteria := util.AqlSearchCriteria{
		Repos:              conf.Repos,
		PropertyFilter:     conf.FilterProperties,
		LastDownloadedFrom: conf.LastDownloadedFrom,
		LastDownloadedTo:   conf.LastDownloadedTo,
	}

	criteriaJson, err := aqlCriteria.GetJson()
	if err != nil {
		return nil, fmt.Errorf("failed to create AQL criteria: %w", err)
	}
	aql := fmt.Sprintf(aqlDownloadTemplate, criteriaJson)

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
	pageSize := conf.PageSize
	workersLock := make(chan bool, conf.MaxConcurrentWorkers)

	var mapperWorkers []*statMapper
	getValueFunc := func(item *util.AqlItem) int {
		return item.Stats[0].Downloads
	}

	scheduledItems := 0
	/*
		While there are new items, filter and map them to the requested identity
	*/
	for scheduledItems < itemsCount {

		// Start mappers in background for each page
		mapper := newStatMapper()
		mapper.GetValueFunc = getValueFunc

		mapperWorkers = append(mapperWorkers, mapper)
		pageInitialIndex := scheduledItems
		pageFinalIndex := scheduledItems + pageSize
		if pageFinalIndex > itemsCount {
			pageFinalIndex = itemsCount
		}

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
