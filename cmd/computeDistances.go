/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/Outer-Rim-Outcast-Coalition/bodies/distances"
)

// computeDistancesCmd represents the computeDistances command
var computeDistancesCmd = &cobra.Command{
	Use:   "computeDistances",
	Short: "Pre-compute system distances to Sol",
	Long: `Reads in a system dump from EDSM (with coords) and writes a gob file
mapping system IDs to distances from Sol.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("computeDistances called")
		in,_ := cmd.Flags().GetString("system-file")
		fmt.Printf("reading from: %s\n", in)
		gob,_ := cmd.Flags().GetString("gob-file")
		fmt.Printf("using gob: %s\n", gob)
		distances.MakeDB(in, gob)
	},
}

func init() {
	rootCmd.AddCommand(computeDistancesCmd)

	var InputFilename string
	var GobFilename string
	computeDistancesCmd.Flags().StringVarP(&InputFilename, "system-file", "s", "", "EDSM gzipped dump to read system data from.")
	computeDistancesCmd.Flags().StringVarP(&GobFilename, "gob-file", "g", "", "Filename to write gob'ed distance data to.")
}
