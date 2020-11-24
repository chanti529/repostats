package commands

import (
	"errors"
	"fmt"
	"github.com/chanti529/jfrog-cli-plugin-template/service"
	"github.com/chanti529/jfrog-cli-plugin-template/util"
	"github.com/cheynewallace/tabby"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"github.com/jfrog/jfrog-cli-core/utils/config"
	"strconv"
	"strings"
	"text/tabwriter"
)

func GetRepoStatDownloadCommand() components.Command {
	return components.Command{
		Name:        "download",
		Description: "Get repo download statistics.",
		Aliases:     []string{"d"},
		Arguments:   getRepoStatDownloadArguments(),
		Flags:       getRepoStatDownloadFlags(),
		//EnvVars:     getHelloEnvVar(),
		Action: func(c *components.Context) error {
			return repoStatDownloadCmd(c)
		},
	}
}

// TODO: Create size command

func getRepoStatDownloadArguments() []components.Argument {
	return []components.Argument{
		{
			Name:        "type",
			Description: "Type of component to get statistics. Valid values: artifact, folder, repo, user",
		},
	}
}

func getRepoStatDownloadFlags() []components.Flag {
	// TODO: Setup additional flags
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
		},
		components.StringFlag{
			Name:         "limit",
			Description:  `Max number or results. Set value to 0 to disable limit`,
			DefaultValue: "5",
		},
		components.StringFlag{
			Name:         "sort",
			Description:  "Results order. Valid values: desc, asc, alpha",
			DefaultValue: "desc",
		},
	}
}

func repoStatDownloadCmd(c *components.Context) error {
	if len(c.Arguments) != 1 {
		return errors.New("Wrong number of arguments. Expected: 1, " + "Received: " + strconv.Itoa(len(c.Arguments)))
	}

	// TODO: Validate arguments and combinations

	// Get Target Artifactory Configuration
	targetRtConfig, err := getTargetArtifactoryConfig(c.GetStringFlagValue("server-id"))
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to get Artifactory configuration: %w", err))
	}

	// TODO: Set command configuration
	conf := service.RepoStatConfiguration{
		RtDetails: targetRtConfig,
		Type:      c.Arguments[0],
		Repos:     strings.Split(c.GetStringFlagValue("repos"), ","),
	}

	limit, err := getIntFlagValue(c, "limit")
	if err != nil {
		return err
	}
	conf.Limit = limit

	// Execute command
	results, err := service.GetDownloadStat(&conf)
	if err != nil {
		return err
	}

	// Write output as a table
	w := tabwriter.NewWriter(&util.LogIoWriter{}, 0, 0, 2, ' ', 0)
	t := tabby.NewCustom(w)
	for _, item := range results {
		t.AddLine(item.Id, item.Value)
	}
	t.Print()
	return nil
}

func getTargetArtifactoryConfig(serverName string) (*config.ArtifactoryDetails, error) {
	return config.GetArtifactorySpecificConfig(serverName, true, false)
}

func getIntFlagValue(c *components.Context, flagName string) (int, error) {
	limit := c.GetStringFlagValue(flagName)
	return strconv.Atoi(limit)
}
