package commands

import (
	"errors"
	"github.com/chanti529/jfrog-cli-plugin-template/service"
	"github.com/chanti529/jfrog-cli-plugin-template/util"
	"github.com/cheynewallace/tabby"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"strconv"
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
	flags := []components.Flag{
		components.StringFlag{
			Name:         "modifiedfrom",
			Description:  "Filter artifacts modified after given timestamp in format RFC3339.",
			DefaultValue: "",
		},
		components.StringFlag{
			Name:         "modifiedto",
			Description:  "Filter artifacts modified before given timestamp in format RFC3339.",
			DefaultValue: "",
		},
	}
	flags = append(flags, getCommonFlags()...)
	return flags
}

func repoStatSizeCmd(c *components.Context) error {
	if len(c.Arguments) != 1 {
		return errors.New("Wrong number of arguments. Expected: 1, " + "Received: " + strconv.Itoa(len(c.Arguments)))
	}

	conf := service.RepoStatConfiguration{
		Type: c.Arguments[0],
	}

	err := parseCommonFlags(c, &conf)
	if err != nil {
		return err
	}

	modifiedFrom, err := getTimestampFlagValue(c, "modifiedfrom")
	if err != nil {
		return err
	}
	conf.ModifiedFrom = modifiedFrom

	modifiedTo, err := getTimestampFlagValue(c, "modifiedto")
	if err != nil {
		return err
	}
	conf.ModifiedTo = modifiedTo

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
