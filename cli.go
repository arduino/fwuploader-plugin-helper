// fwuploader-plugin-helper
// Copyright (c) 2023 Arduino LLC.  All right reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package helper

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/arduino/go-paths-helper"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
	"gopkg.in/yaml.v3"
)

// RunPlugin runs the given plugin
func RunPlugin(plugin Plugin) {
	info := plugin.GetPluginInfo()

	var (
		fqbn        string
		logLevel    string
		verbose     bool
		portAddress string
	)

	firmwareFlashCmd := &cobra.Command{
		Use:   "flash",
		Short: "Upload a firmware to the board",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fatal("Please specify firmware path", 2)
			}
			fwPath := paths.New(args[0])
			err := plugin.UploadFirmware(
				portAddress,
				fqbn,
				fwPath,
				&PluginFeedback{stdOut: os.Stdout, stdErr: os.Stderr},
			)
			if err != nil {
				fatal(err.Error(), 3)
			}
		},
	}

	firmwareGetVersionCmd := &cobra.Command{
		Use:   "get-version",
		Short: "Get the version of the currently installed firmware on the board",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				fatal("Invalid arguments provided", 4)
			}
			version, err := plugin.GetFirmwareVersion(
				portAddress,
				fqbn,
				&PluginFeedback{stdOut: os.Stdout, stdErr: os.Stderr},
			)
			if err != nil {
				fatal(err.Error(), 3)
			}
			fmt.Println("FIRMWARE-VERSION:", version)
		},
	}

	firmwareCmd := &cobra.Command{
		Use:   "firmware",
		Short: "Firmware handling commands",
	}
	firmwareCmd.AddCommand(firmwareFlashCmd)
	firmwareCmd.AddCommand(firmwareGetVersionCmd)

	certFlashCmd := &cobra.Command{
		Use:   "flash",
		Short: "Upload a certificate on the board",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fatal("Please specify a certificate path", 2)
			}
			certPath := paths.New(args[0])
			err := plugin.UploadCertificate(
				portAddress,
				fqbn,
				certPath,
				&PluginFeedback{stdOut: os.Stdout, stdErr: os.Stderr},
			)
			if err != nil {
				fatal(err.Error(), 3)
			}
		},
	}

	certCmd := &cobra.Command{
		Use:   "cert",
		Short: "Certificates handling commands",
	}
	certCmd.AddCommand(certFlashCmd)

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Return informations about this fw-updater plugin",
		Run: func(cmd *cobra.Command, args []string) {
			printInfo(info)
		},
	}

	appName := filepath.Base(os.Args[0])
	cli := &cobra.Command{
		Use:   appName,
		Short: info.Name + " - This is an Arduino Firmware Uploader plugin.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			t, found := map[string]slog.Level{
				"trace": slog.LevelDebug,
				"debug": slog.LevelDebug,
				"info":  slog.LevelInfo,
				"warn":  slog.LevelWarn,
				"error": slog.LevelError,
				"fatal": slog.LevelError,
				"panic": slog.LevelError,
			}[logLevel]
			if !found {
				t = slog.LevelError
			}

			var w io.Writer
			if !verbose {
				w = io.Discard
			} else {
				w = os.Stdout
			}
			slog.SetDefault(slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{Level: t})))
		},
	}
	cli.AddCommand(firmwareCmd)
	cli.AddCommand(certCmd)
	cli.AddCommand(versionCmd)
	cli.PersistentFlags().StringVarP(&portAddress, "address", "p", "", "Port address")
	// The fqbn is an optional flag that can be used by the plugin to do specific operations with a board.
	// With this input we can support a family of boards and not only a single one per plugin
	cli.PersistentFlags().StringVarP(&fqbn, "fqbn", "b", "", "Fully qualified board name")
	cli.PersistentFlags().StringVar(&logLevel, "log-level", "info", "Messages with this level and above will be logged. Valid levels are: trace, debug, info, warn, error, fatal, panic")
	cli.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Print the logs on the standard output.")

	if err := cli.Execute(); err != nil {
		fatal(err.Error(), 1)
	}
}

func fatal(msg string, exitcode int) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", msg)
	os.Exit(exitcode)
}

func printInfo(info *PluginInfo) {
	type infoResult struct {
		PluginInfo       *PluginInfo `yaml:"plugin_info"`
		PluginAPIVersion int         `yaml:"plugin_api_version"`
	}
	data, err := yaml.Marshal(&infoResult{
		PluginAPIVersion: 1,
		PluginInfo:       info,
	})
	if err != nil {
		fatal(err.Error(), 3)
	}
	fmt.Println(string(data))
}
