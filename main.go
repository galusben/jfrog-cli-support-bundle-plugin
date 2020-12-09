package main

import (
	"github.com/jfrog/jfrog-cli-core/plugins"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"github.com/jfrog/jfrog-cli-plugin-template/commands"
)

func main() {
	plugins.PluginMain(getApp())
}

func getApp() components.App {
	app := components.App{}
	app.Name = "support"
	app.Description = "Perform support operations like creating and uploading support bundles"
	app.Version = "v0.0.1"
	app.Commands = getCommands()
	return app
}

func getCommands() []components.Command {
	return []components.Command{
		commands.GetUploadSupportBundleCommand(),
		commands.GetGenerateSupportBundleCommand(),
		commands.GetDecryptCommand(),
		commands.GetEncryptCommand(),
	}
}
