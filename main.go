package main

import (
	"github.com/chanti529/repostats/commands"
	"github.com/jfrog/jfrog-cli-core/v2/plugins"
	"github.com/jfrog/jfrog-cli-core/v2/plugins/components"
)

func main() {
	plugins.PluginMain(getApp())
}

func getApp() components.App {
	app := components.App{}
	app.Name = "repostats"
	app.Description = "Get Artifacts statistics."
	app.Version = "v1.0.3"
	app.Commands = getCommands()
	return app
}

func getCommands() []components.Command {
	return []components.Command{
		commands.GetRepoStatDownloadCommand(),
		commands.GetRepoStatSizeCommand(),
	}
}
