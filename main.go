package main

import (
	"fmt"
	"os"

	prompt "github.com/c-bata/go-prompt"
	cli "github.com/jawher/mow.cli"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"

	"github.com/InjectiveLabs/dexterm/version"
)

var app = cli.App("dexterm", "A command line client to InjectiveProtocol DEX")

func main() {
	app.Action = func() {
		if toBool(*appConfigMap["log.debug"]) {
			logrus.SetLevel(logrus.TraceLevel)
		}

		configPath, _ := homedir.Expand(*configPath)
		ctl, err := NewAppController(configPath)
		if err != nil {
			logrus.Fatalln(err)
		}

		state := NewAppState(ctl)

		fmt.Println("Welcome to DEXTerm! Use tab to autocomplete commands. Ctrl-D to quit.")
		startPrompt(state)
	}

	app.Command("v version", "Print application version", versionCmd)

	if err := app.Run(os.Args); err != nil {
		logrus.Fatalln(err)
	}
}

func versionCmd(c *cli.Cmd) {
	c.Action = func() {
		fmt.Println(version.Version())
	}
}

func startPrompt(state *AppState) {
	p := prompt.New(
		state.Executor(),
		state.Completer(),
		prompt.OptionLivePrefix(state.LivePrefix()),
		prompt.OptionTitle("DEXTerm"),
		prompt.OptionPrefix("∆ "),
		prompt.OptionShowCompletionAtStart(),
		prompt.OptionPrefixTextColor(prompt.DarkBlue),
		prompt.OptionSuggestionBGColor(prompt.DarkBlue),
		prompt.OptionAddKeyBind(prompt.KeyBind{
			Key: prompt.ControlC,
			Fn: func(buf *prompt.Buffer) {
				state.DiscardCmd()
			},
		}),
	)

	p.Run()
}
