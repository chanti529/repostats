package service

import (
	"errors"
	"github.com/chanti529/jfrog-cli-plugin-template/util"
	"github.com/jfrog/jfrog-cli-core/utils/config"
	"time"
)

type StatItem struct {
	Id    string
	Value int
}

type RepoStatConfiguration struct {
	RtDetails *config.ArtifactoryDetails
	Type      string
	Repos     []string
	Limit     int
}

type statMapper struct {
	GetValueFunc func(item *util.AqlItem) int
	Done         bool
	Error        error
	Result       map[string]int
}

func (w *statMapper) process(items []*util.AqlItem, conf *RepoStatConfiguration) {
	if w.GetValueFunc == nil {
		w.Error = errors.New("Get Value function not set on mapper")
		return
	}

	for _, item := range items {

		// TODO: Filter item

		itemId, err := getItemIdentity(item, conf)
		if err != nil {
			w.Error = err
		}

		value := w.GetValueFunc(item)

		// Map result
		w.Result[itemId] = w.Result[itemId] + value
	}

	// TODO: Apply sort and limit if type is artifact since it cannot be reduced any further

	w.Done = true
}

func waitForWorkers(workers []*statMapper) error {
	for _, worker := range workers {
		for !worker.Done && worker.Error == nil {
			time.Sleep(time.Second)
		}
		if worker.Error != nil {
			return worker.Error
		}
	}

	return nil
}

func getItemIdentity(item *util.AqlItem, conf *RepoStatConfiguration) (string, error) {
	switch conf.Type {
	case "repo":
		return item.Repo, nil
	case "folder":
		//TODO: Get folder identity
		return "", errors.New("Not implemented")
	case "artifact":
		return item.GetFullPath(), nil
	case "user":
		return item.ModifiedBy, nil
	default:
		return "", errors.New("Invalid type")
	}
}

func reduce(workers []*statMapper) ([]StatItem, error) {
	mergedResults := workers[0].Result

	for workerIndex := 1; workerIndex < len(workers); workerIndex++ {
		worker := workers[workerIndex]

		// Merge results between workers
		for id, value := range worker.Result {
			mergedResults[id] = mergedResults[id] + value
		}
	}

	// TODO: Apply sort and limit
	statItems := []StatItem{}
	for id, value := range mergedResults {
		statItems = append(statItems, StatItem{
			Id:    id,
			Value: value,
		})
	}

	return statItems, nil
}
