/*
Copyright AppsCode Inc. and Contributors

Licensed under the AppsCode Community License 1.0.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://github.com/appscode/licenses/raw/1.0.0/AppsCode-Community-1.0.0.md

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmds

import (
	"flag"
	"os"

	"kubedb.dev/apimachinery/client/clientset/versioned/scheme"

	"github.com/appscode/go/flags"
	"github.com/appscode/go/log"
	"github.com/appscode/go/log/golog"
	v "github.com/appscode/go/version"
	"github.com/spf13/cobra"
	genericapiserver "k8s.io/apiserver/pkg/server"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"
	"kmodules.xyz/client-go/logs"
	"kmodules.xyz/client-go/tools/cli"
	appcatscheme "kmodules.xyz/custom-resources/client/clientset/versioned/scheme"
)

func NewRootCmd(version string) *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:               "percona-xtradb-operator",
		DisableAutoGenTag: true,
		PersistentPreRun: func(c *cobra.Command, args []string) {
			flags.DumpAll(c.Flags())
			cli.SendAnalytics(c, version)

			//scheme.AddToScheme(clientsetscheme.Scheme)
			if err := scheme.AddToScheme(clientsetscheme.Scheme); err != nil {
				log.Errorln(err)
			}
			//appcatscheme.AddToScheme(clientsetscheme.Scheme)
			if err := appcatscheme.AddToScheme(clientsetscheme.Scheme); err != nil {
				log.Errorln(err)
			}
			cli.LoggerOptions = golog.ParseFlags(c.Flags())
		},
	}
	rootCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
	// ref: https://github.com/kubernetes/kubernetes/issues/17162#issuecomment-225596212
	logs.ParseFlags()
	rootCmd.PersistentFlags().BoolVar(&cli.EnableAnalytics, "enable-analytics", cli.EnableAnalytics, "Send analytical events to Google Analytics")

	rootCmd.AddCommand(v.NewCmdVersion())

	stopCh := genericapiserver.SetupSignalHandler()
	rootCmd.AddCommand(NewCmdRun(version, os.Stdout, os.Stderr, stopCh))

	return rootCmd
}
