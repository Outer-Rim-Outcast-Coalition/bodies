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
	"github.com/Outer-Rim-Outcast-Coalition/bodies/filter"
)

// filterBodiesCmd represents the filterBodies command
var filterBodiesCmd = &cobra.Command{
	Use:   "filterBodies",
	Short: "Filter bodies for candidates to have crystalline shards.",
	Long: `Reads in a bodies dump from EDSM and gob file from computeDistances,
then filters the bodies in the dump according to the criteria for the existence
of Crystalline Shards. It outputs a sorted list in either JSON (all fields) or
CSV (selected fields). Can also be used to re-export the JSON to CSV.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("filterBodies called")
		in,_ := cmd.Flags().GetString("bodies-file")
		fmt.Printf("reading from: %s\n", in)
		gob,_ := cmd.Flags().GetString("gob-file")
		fmt.Printf("using gob: %s\n", gob)
		out,_ := cmd.Flags().GetString("output")
		fmt.Printf("using output: %s\n", out)
		oft,_ := cmd.Flags().GetString("format")
		fmt.Printf("using format: %s\n", oft)
		rex,_ := cmd.Flags().GetString("reexport")
		fmt.Printf("using re-export: %s\n", rex)
		lim,_ := cmd.Flags().GetInt64("limit")
		fmt.Printf("using limit: %d\n", lim)
		filter.FilterBodies(in, gob, out, oft, rex, lim)
	},
}

func init() {
	rootCmd.AddCommand(filterBodiesCmd)

	var InputFilename string
	var GobFilename string
	var OutputFilename string
	var OutputFormat string
	var ReExport string
	var BodyLimit int64
	filterBodiesCmd.Flags().StringVarP(&InputFilename, "bodies-file", "b", "", "EDSM gzipped dump to read body system data from.")
	filterBodiesCmd.Flags().StringVarP(&GobFilename, "gob-file", "g", "", "Filename to write gob'ed distance data to.")
	filterBodiesCmd.Flags().StringVarP(&OutputFilename, "output", "o", "", "Filename to write data to.")
	filterBodiesCmd.Flags().StringVarP(&OutputFormat, "format", "f", "", "Output data format.")
	filterBodiesCmd.Flags().StringVarP(&ReExport, "reexport", "r", "", "Candidates JSON dump to re-export to new format.")
	filterBodiesCmd.Flags().Int64VarP(&BodyLimit, "limit", "l", 0, "Limit number of bodies to process.")
}
