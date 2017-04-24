package main

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/joho/godotenv"
	"github.com/urfave/cli"
)

var build = "0" // build number set at compile-time

func main() {
	app := cli.NewApp()
	app.Name = "helm plugin"
	app.Usage = "helm plugin"
	app.Action = run
	app.Version = fmt.Sprintf("1.0.%s", build)
	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:   "helm_command",
			Usage:  "add the command Helm has to execute",
			EnvVar: "PLUGIN_HELM_COMMAND,HELM_COMMAND",
		},
		cli.StringFlag{
			Name:   "namespace",
			Usage:  "Kubernetes namespace",
			EnvVar: "PLUGIN_NAMESPACE,NAMESPACE",
		},
		cli.StringFlag{
			Name:   "release",
			Usage:  "Kubernetes helm release",
			EnvVar: "PLUGIN_RELEASE,RELEASE",
		},
		cli.StringFlag{
			Name:   "chart",
			Usage:  "Kubernetes helm release",
			EnvVar: "PLUGIN_CHART,CHART",
		},
		cli.StringFlag{
			Name:   "values",
			Usage:  "Kubernetes helm release",
			EnvVar: "PLUGIN_VALUES,VALUES",
		},
		cli.StringFlag{
			Name:   "values_files",
			Usage:  "Helm values override files",
			EnvVar: "PLUGIN_VALUES_FILES,VALUES_FILES",
		},
		cli.BoolFlag{
			Name:   "skip_tls_verify",
			Usage:  "Skip TLS verification",
			EnvVar: "PLUGIN_SKIP_TLS_VERIFY,SKIP_TLS_VERIFY",
		},
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "Debug",
			EnvVar: "PLUGIN_DEBUG,DEBUG",
		},
		cli.BoolFlag{
			Name:   "dry-run",
			Usage:  "Helm dry-run",
			EnvVar: "PLUGIN_DRY_RUN,DRY_RUN",
		},
		cli.StringFlag{
			Name:   "prefix",
			Usage:  "Prefix for all the secrets",
			EnvVar: "PLUGIN_PREFIX,PREFIX",
		},
		cli.StringFlag{
			Name:   "tiller-ns",
			Usage:  "Namespace to install Tiller",
			EnvVar: "PLUGIN_TILLER_NS,TILLER_NS",
		},
		cli.BoolFlag{
			Name:   "wait",
			Usage:  "if set, will wait until all Pods, PVCs, and Services are in a ready state before marking the release as successful.",
			EnvVar: "PLUGIN_WAIT,WAIT",
		},
	}
	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(c *cli.Context) error {
	if c.String("env-file") != "" {
		_ = godotenv.Load(c.String("env-file"))
	}
	plugin := Plugin{
		Config: Config{
			APIServer:     c.String("api_server"),
			Token:         c.String("token"),
			HelmCommand:   c.StringSlice("helm_command"),
			Namespace:     c.String("namespace"),
			SkipTLSVerify: c.Bool("skip_tls_verify"),
			Values:        c.String("values"),
			ValuesFiles:   c.String("values_files"),
			Release:       c.String("release"),
			Chart:         c.String("chart"),
			Debug:         c.Bool("debug"),
			DryRun:        c.Bool("dry-run"),
			Secrets:       c.StringSlice("secrets"),
			Prefix:        c.String("prefix"),
			TillerNs:      c.String("tiller-ns"),
			Wait:          c.Bool("wait"),
		},
	}
	resolveSecrets(&plugin)
	if plugin.Config.Debug {
		plugin.debug()
	}
	return plugin.Exec()
}
