package commands

import (
	"errors"
	"strconv"
	"text/tabwriter"

	"github.com/chanti529/repostats/service"
	"github.com/chanti529/repostats/util"
	"github.com/cheynewallace/tabby"
	"github.com/jfrog/jfrog-cli-core/v2/artifactory/utils"
	"github.com/jfrog/jfrog-cli-core/v2/plugins/components"
)

func GetRepoStatDownloadCommand() components.Command {
	return components.Command{
		Name:        "download",
		Description: "Get repo download count statistics.",
		Aliases:     []string{"d"},
		Arguments:   getRepoStatDownloadArguments(),
		Flags:       getRepoStatDownloadFlags(),
		Action: func(c *components.Context) error {
			return repoStatDownloadCmd(c)
		},
	}
}

func getRepoStatDownloadArguments() []components.Argument {
	return []components.Argument{
		{
			Name:        "type",
			Description: "Type of component to get statistics. Valid values: artifact, folder, repo, user",
		},
	}
}

func getRepoStatDownloadFlags() []components.Flag {
	flags := []components.Flag{
		components.StringFlag{
			Name:         "lastdownloadedfrom",
			Description:  "Filter artifacts last downloaded after given timestamp in RFC3339 format.",
			DefaultValue: "",
		},
		components.StringFlag{
			Name:         "lastdownloadedto",
			Description:  "Filter artifacts last downloaded before given timestamp in RFC3339 format.",
			DefaultValue: "",
		},
	}
	flags = append(flags, getCommonFlags()...)
	return flags
}

func repoStatDownloadCmd(c *components.Context) error {
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

	lastDownloadedFrom, err := getTimestampFlagValue(c, "lastdownloadedfrom")
	if err != nil {
		return err
	}
	conf.LastDownloadedFrom = lastDownloadedFrom

	lastDownloadedTo, err := getTimestampFlagValue(c, "lastdownloadedto")
	if err != nil {
		return err
	}
	conf.LastDownloadedTo = lastDownloadedTo

	servicesManager, err := utils.CreateServiceManager(conf.RtDetails, 5, 200, false)
	if err != nil {
		return err
	}

	results, err := service.GetDownloadStat(&conf, servicesManager)
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
