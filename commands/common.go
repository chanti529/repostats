package commands

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/chanti529/repostats/service"
	"github.com/chanti529/repostats/util"
	"github.com/jfrog/jfrog-cli-core/v2/plugins/components"
	"github.com/jfrog/jfrog-cli-core/v2/utils/config"
)

func getCommonFlags() []components.Flag {
	return []components.Flag{
		components.StringFlag{
			Name:         "server-id",
			Description:  "Artifactory server ID configured using the config command.",
			DefaultValue: "",
		},
		components.StringFlag{
			Name:         "repos",
			Description:  "Comma separated list of repositories.",
			DefaultValue: "",
			Mandatory:    true,
		},
		components.StringFlag{
			Name:         "path",
			Description:  "Regular Expression to filter the full path of artifacts.",
			DefaultValue: "",
		},
		components.StringFlag{
			Name:         "properties",
			Description:  "Comma separeted list of properties and values to filter in the format property_name=pattern",
			DefaultValue: "",
		},
		components.StringFlag{
			Name:         "limit",
			Description:  `Max number or results. Set value to 0 to disable limit`,
			DefaultValue: "10",
		},
		components.StringFlag{
			Name:         "sort",
			Description:  "Results order. Valid values: desc, asc, alpha",
			DefaultValue: "desc",
		},
		components.StringFlag{
			Name:         "page-size",
			Description:  "Number of items to be processed at once per a single worker",
			DefaultValue: "50000",
		},
		components.StringFlag{
			Name:         "max-workers",
			Description:  "Max number of concurrent workers processing items in parallel at a given time",
			DefaultValue: "5",
		},
		components.StringFlag{
			Name:         "max-depth",
			Description:  "Max depth to group folders when using folder command type",
			DefaultValue: "4",
		},
	}
}

func parseCommonFlags(c *components.Context, conf *service.RepoStatConfiguration) error {
	conf.Repos = strings.Split(c.GetStringFlagValue("repos"), ",")
	conf.Sort = c.GetStringFlagValue("sort")

	// Get Target Artifactory Configuration
	targetRtConfig, err := getTargetArtifactoryConfig(c.GetStringFlagValue("server-id"))
	if err != nil {
		return fmt.Errorf("Failed to get Artifactory configuration: %w", err)
	}
	conf.RtDetails = targetRtConfig

	// Parse limit
	limit, err := getIntFlagValue(c, "limit")
	if err != nil {
		return err
	}
	conf.Limit = limit

	// Parse path filter
	filterPath := c.GetStringFlagValue("path")
	if filterPath != "" {
		conf.FilterPathRegexp = regexp.MustCompile(filterPath)
	}

	// Parse properties filter
	filterProperties := c.GetStringFlagValue("properties")
	if filterProperties != "" {
		var filterPropertiesKeyValuePairs []util.KeyValuePair
		filterPropertiesParts := strings.Split(filterProperties, ",")
		for _, item := range filterPropertiesParts {
			filterPropertyItemParts := strings.SplitN(item, "=", 2)
			filterPropertiesKeyValuePairs = append(filterPropertiesKeyValuePairs, util.KeyValuePair{
				Key:   filterPropertyItemParts[0],
				Value: filterPropertyItemParts[1],
			})
		}
		conf.FilterProperties = filterPropertiesKeyValuePairs
	}

	// Parse page size and workers limit
	pageSize, err := getIntFlagValue(c, "page-size")
	if err != nil {
		return err
	}
	conf.PageSize = pageSize

	maxWorkers, err := getIntFlagValue(c, "max-workers")
	if err != nil {
		return err
	}
	conf.MaxConcurrentWorkers = maxWorkers

	maxDepth, err := getIntFlagValue(c, "max-depth")
	if err != nil {
		return err
	}
	conf.MaxDepth = maxDepth

	return nil
}

func getTargetArtifactoryConfig(serverName string) (*config.ServerDetails, error) {
	return config.GetSpecificConfig(serverName, true, false)
}

func getIntFlagValue(c *components.Context, flagName string) (int, error) {
	limit := c.GetStringFlagValue(flagName)
	return strconv.Atoi(limit)
}

func getTimestampFlagValue(c *components.Context, flagName string) (timeValue time.Time, err error) {
	strValue := c.GetStringFlagValue(flagName)
	if strValue != "" {
		timeValue, err = time.Parse(time.RFC3339, strValue)
	}
	return
}
