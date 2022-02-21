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

// combineCmd represents the combine command
var combineCmd = &cobra.Command{
	Use:   "combine",
	Short: "Combine parts to reconstruct original input",
	Long: `Combine a threshold number of parts to reconstruct
original input.`,
	RunE: run.Combine,
}

func init() {
	rootCmd.AddCommand(combineCmd)
	f := combineCmd.Flags()

	f.Int(flags.NumThreshold, 3, "Number of threshold (min=2)")
}
