package commands

import (
	"errors"
	"fmt"
	"github.com/chanti529/jfrog-cli-plugin-template/service"
	"github.com/chanti529/jfrog-cli-plugin-template/util"
	"github.com/cheynewallace/tabby"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"strconv"
	"strings"
	"text/tabwriter"
)

func GetRepoStatSizeCommand() components.Command {
	return components.Command{
		Name:        "size",
		Description: "Get repo size statistics.",
		Aliases:     []string{"s"},
		Arguments:   getRepoStatSizeArguments(),
		Flags:       getRepoStatSizeFlags(),
		//EnvVars:     getHelloEnvVar(),
		Action: func(c *components.Context) error {
			return repoStatSizeCmd(c)
		},
	}
}

func getRepoStatSizeArguments() []components.Argument {
	return []components.Argument{
		{
			Name:        "type",
			Description: "Type of component to get statistics. Valid values: artifact, folder, repo, user",
		},
	}
}

func getRepoStatSizeFlags() []components.Flag {
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
			DefaultValue: "10",
		},
		components.StringFlag{
			Name:         "sort",
			Description:  "Results order. Valid values: desc, asc, alpha",
			DefaultValue: "desc",
		},
	}
}

func repoStatSizeCmd(c *components.Context) error {
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
		Sort:      c.GetStringFlagValue("sort"),
	}

	limit, err := getIntFlagValue(c, "limit")
	if err != nil {
		return err
	}
	conf.Limit = limit

	// Execute command
	results, err := service.GetSizeStat(&conf)
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
