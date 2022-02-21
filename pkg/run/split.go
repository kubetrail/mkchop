package run

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/shamir"
	"github.com/kubetrail/mkchop/pkg/flags"
	"github.com/mr-tron/base58"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Split(cmd *cobra.Command, args []string) error {
	_ = viper.BindPFlag(flags.NumSplits, cmd.Flags().Lookup(flags.NumSplits))
	_ = viper.BindPFlag(flags.NumThreshold, cmd.Flags().Lookup(flags.NumThreshold))

	numSplits := viper.GetInt(flags.NumSplits)
	numThreshold := viper.GetInt(flags.NumThreshold)

	if numSplits < 3 {
		return fmt.Errorf("num-splits needs to be at least 3")
	}

	if numThreshold < 2 {
		return fmt.Errorf("num-threshold needs to be at least 2")
	}

	if numThreshold >= numSplits {
		return fmt.Errorf("num-splits needs to be greater than num-threshold")
	}

	if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Enter input string: "); err != nil {
		return fmt.Errorf("failed to write to output: %w", err)
	}

	inputReader := bufio.NewReader(cmd.InOrStdin())
	input, err := inputReader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read from input: %w", err)
	}
	input = strings.Trim(input, "\n")

	parts, err := shamir.Split([]byte(input), numSplits, numThreshold)
	if err != nil {
		return fmt.Errorf("failed to split input: %w", err)
	}

	for _, part := range parts {
		if _, err := fmt.Fprintln(cmd.OutOrStdout(), base58.Encode(part)); err != nil {
			return fmt.Errorf("failed to write to output: %w", err)
		}
	}

	return nil
}
