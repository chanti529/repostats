package service

import (
	"encoding/json"
	"fmt"
	"github.com/chanti529/repostats/util"
	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"io/ioutil"
)

const (
	aqlSizeTemplate = `items.find(%s).include("repo", "path", "name", "created", "modified", "modified_by", "size").sort({
			"$asc":["created"]
		})`
)

func GetSizeStat(conf *RepoStatConfiguration, servicesManager artifactory.ArtifactoryServicesManager) ([]StatItem, error) {
	aqlCriteria := util.AqlSearchCriteria{
		Repos:          conf.Repos,
		PropertyFilter: conf.FilterProperties,
		ModifiedFrom:   conf.ModifiedFrom,
		ModifiedTo:     conf.ModifiedTo,
	}

	criteriaJson, err := aqlCriteria.GetJson()
	if err != nil {
		return nil, fmt.Errorf("failed to create AQL criteria: %w", err)
	}
	aql := fmt.Sprintf(aqlSizeTemplate, criteriaJson)

	pageSize := conf.PageSize
	workersLock := make(chan bool, conf.MaxConcurrentWorkers)

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

			// Start mappers in background
			mapper := newStatMapper()
			mapper.GetValueFunc = getValueFunc
			mapperWorkers = append(mapperWorkers, mapper)

			workersLock <- true
			go func(m *statMapper, items []*util.AqlItem) {
				m.process(items, conf)
				<-workersLock
			}(mapper, parsedResult.Results)
		}

		log.Debug(fmt.Sprintf("Found %v artifacts...", itemsCount))
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
