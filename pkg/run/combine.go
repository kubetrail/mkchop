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

func Combine(cmd *cobra.Command, args []string) error {
	_ = viper.BindPFlag(flags.NumThreshold, cmd.Flags().Lookup(flags.NumThreshold))
	numThreshold := viper.GetInt(flags.NumThreshold)

	if numThreshold < 2 {
		return fmt.Errorf("num-threshold has to be at least 2")
	}

	inputReader := bufio.NewReader(cmd.InOrStdin())
	partsMap := make(map[string]struct{})
	parts := make([][]byte, numThreshold)
	for i := 0; i < numThreshold; i++ {
		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Enter part %d of %d: ", i+1, numThreshold); err != nil {
			return fmt.Errorf("failed to write to output: %w", err)
		}
		input, err := inputReader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read from input: %w", err)
		}
		input = strings.Trim(input, "\n")
		partsMap[input] = struct{}{}
		parts[i], err = base58.Decode(input)
		if err != nil {
			return fmt.Errorf("failed to base58 decode part %d: %w", i, err)
		}
	}

	if len(partsMap) != numThreshold {
		return fmt.Errorf("please input at least %d unique parts", numThreshold)
	}

	out, err := shamir.Combine(parts)
	if err != nil {
		return fmt.Errorf("failed to combine parts: %w", err)
	}

	if _, err := fmt.Fprintln(cmd.OutOrStdout(), string(out)); err != nil {
		return fmt.Errorf("failed to write to output: %w", err)
	}

	return nil
}
