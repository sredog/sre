/*
Copyright © 2022 Mikolaj Pawlikowski <mikolaj@pawlikowski.pl>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"

	"github.com/prometheus/procfs"
	"github.com/spf13/cobra"
	"github.com/sredog/sre/pkg/analysis"
	"github.com/sredog/sre/pkg/loadavg"
	"github.com/sredog/sre/pkg/uptime"
)

// quickCmd represents the quick command
var quickCmd = &cobra.Command{
	Use:   "quick",
	Short: "Quick overview of the system: CPUs, RAM, IO, net, filesystems",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		p, err := procfs.NewDefaultFS()
		if err != nil {
			panic(err)
		}

		var probes []analysis.Probe

		uptime, err := uptime.NewUptime(2)
		if err != nil {
			panic(err)
		}
		probes = append(probes, uptime)

		la, err := loadavg.NewLoadAverage(p)
		if err != nil {
			panic(err)
		}
		probes = append(probes, la)

		for _, probe := range probes {
			output := probe.Display()
			_, err = fmt.Print(output)
			if err != nil {
				panic(err)
			}
			for _, observation := range probe.Analysis() {
				_, err = fmt.Printf("%s\n", observation.Format())
				if err != nil {
					panic(err)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(quickCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// quickCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// quickCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
