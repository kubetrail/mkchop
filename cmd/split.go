/*
Copyright © 2022 kubetrail.io authors

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
	"github.com/kubetrail/mkchop/pkg/flags"
	"github.com/kubetrail/mkchop/pkg/run"
	"github.com/spf13/cobra"
)

// splitCmd represents the split command
var splitCmd = &cobra.Command{
	Use:   "split",
	Short: "Split input using Shamir splitting algo",
	Long: `
Shamir split algorithm is a way to split input data into a number
of parts such that a minimum threshold of those can be used to
regenerate the original data.

For instance, if input is split into 5 parts with a threshold of
3, it would only require any three parts to reconstruct the original
data.
`,
	RunE: run.Split,
}

func init() {
	rootCmd.AddCommand(splitCmd)
	f := splitCmd.Flags()

	f.Int(flags.NumSplits, 5, "Number of splits")
	f.Int(flags.NumThreshold, 3, "Number of threshold (has to be less than num-splits, min=2)")
	f.String(flags.Key, "", "Input key to split")
}
