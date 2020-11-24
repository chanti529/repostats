package service

import (
	"errors"
	"github.com/chanti529/jfrog-cli-plugin-template/util"
	"github.com/jfrog/jfrog-cli-core/utils/config"
	"time"
)

const (
	TypeArtifact = "artifact"
	TypeFolder   = "folder"
	TypeRepo     = "repo"
	TypeUser     = "user"

	SortAlpha = "alpha"
	SortAsc   = "asc"
	SortDesc  = "desc"
)

type StatItem struct {
	Id    string
	Value int
}

type RepoStatConfiguration struct {
	RtDetails *config.ArtifactoryDetails
	Type      string
	Repos     []string
	Sort      string
	Limit     int
}

type statMapper struct {
	GetValueFunc func(item *util.AqlItem) int
	Done         bool
	Error        error
	Result       map[string]int
}

func newStatMapper() *statMapper {
	return &statMapper{
		Result: make(map[string]int),
	}
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

	// Apply sort and limit if type is artifact since it cannot be reduced any further
	if conf.Type == TypeArtifact && conf.Limit > 0 && conf.Limit < len(w.Result) {

		workerResults := sortAndLimit(conf.Sort, conf.Limit, w.Result)
		// Reset worker result
		w.Result = make(map[string]int)
		for _, item := range workerResults {
			w.Result[item.Id] = item.Value
		}
	}

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

func reduce(workers []*statMapper, conf *RepoStatConfiguration) ([]StatItem, error) {
	mergedResults := workers[0].Result

	// Merge workers results
	for workerIndex := 1; workerIndex < len(workers); workerIndex++ {
		worker := workers[workerIndex]

		// Merge results between workers
		for id, value := range worker.Result {
			mergedResults[id] = mergedResults[id] + value
		}

		// Free worker result for GC
		worker.Result = nil
	}

	var resultsSize int
	if conf.Limit > 0 {
		resultsSize = conf.Limit
	} else {
		resultsSize = len(mergedResults)
	}

	results := sortAndLimit(conf.Sort, resultsSize, mergedResults)
	return results, nil
}

func sortAndLimit(sort string, limit int, items map[string]int) []StatItem {
	results := make([]StatItem, limit)

	// Insert items at their right position in result slice
	for id, value := range items {

		var itemIndex = -1

		// Find item position in results
		for resultIndex, resultItem := range results {

			// In case result position is free use it
			if resultItem.Id == "" {
				itemIndex = resultIndex
				break
			}

			// Descending Sort
			if (sort == SortDesc && value >= resultItem.Value) ||
				// Ascending Sort
				(sort == SortAsc && value <= resultItem.Value) ||
				// Alphabetic Sort
				(sort == SortAlpha && id < resultItem.Id) {
				itemIndex = resultIndex
				break
			}
		}

		// If item has a position in results
		if itemIndex > -1 {

			// If item position is not last we need to shift the results at its right
			if itemIndex != limit-1 {
				copy(results[itemIndex+1:limit], results[itemIndex:limit-1])
			}

			// Add item to results
			results[itemIndex] = StatItem{
				Id:    id,
				Value: value,
			}
		}
	}

	return results
}
