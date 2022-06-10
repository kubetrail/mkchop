package run

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/shamir"
	"github.com/kubetrail/bip39/pkg/mnemonics"
	"github.com/kubetrail/bip39/pkg/prompts"
	"github.com/kubetrail/mkchop/pkg/flags"
	"github.com/mr-tron/base58"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func Split(cmd *cobra.Command, args []string) error {
	persistentFlags := getPersistentFlags(cmd)

	_ = viper.BindPFlag(flags.NumSplits, cmd.Flag(flags.NumSplits))
	_ = viper.BindPFlag(flags.NumThreshold, cmd.Flag(flags.NumThreshold))
	_ = viper.BindPFlag(flags.Key, cmd.Flag(flags.Key))

	numSplits := viper.GetInt(flags.NumSplits)
	numThreshold := viper.GetInt(flags.NumThreshold)
	key := viper.GetString(flags.Key)

	if numSplits < 3 {
		return fmt.Errorf("num-splits needs to be at least 3")
	}

	if numThreshold < 2 {
		return fmt.Errorf("num-threshold needs to be at least 2")
	}

	if numThreshold >= numSplits {
		return fmt.Errorf("num-splits needs to be greater than num-threshold")
	}

	prompt, err := prompts.Status()
	if err != nil {
		return fmt.Errorf("failed to get prompt status: %w", err)
	}

	if len(key) == 0 {
		if len(args) == 0 {
			if prompt {
				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Enter key to split: "); err != nil {
					return fmt.Errorf("failed to write to output: %w", err)
				}
			}
			key, err = mnemonics.Read(cmd.InOrStdin())
			if err != nil {
				return fmt.Errorf("failed to read pub key from input: %w", err)
			}
		} else {
			key = mnemonics.NewFromFields(args)
		}
	}

	parts, err := shamir.Split([]byte(key), numSplits, numThreshold)
	if err != nil {
		return fmt.Errorf("failed to split input: %w", err)
	}

	type output struct {
		Parts []string `json:"parts,omitempty" yaml:"parts,omitempty"`
	}

	out := &output{Parts: make([]string, len(parts))}

	for i, part := range parts {
		out.Parts[i] = base58.Encode(part)
	}

	switch strings.ToLower(persistentFlags.OutputFormat) {
	case flags.OutputFormatNative:
		for i := range out.Parts {
			if _, err := fmt.Fprintln(cmd.OutOrStdout(), out.Parts[i]); err != nil {
				return fmt.Errorf("failed to write part to output: %w", err)
			}
		}
	case flags.OutputFormatYaml:
		jb, err := yaml.Marshal(out)
		if err != nil {
			return fmt.Errorf("failed to serialize output to yaml: %w", err)
		}

		if _, err := fmt.Fprint(cmd.OutOrStdout(), string(jb)); err != nil {
			return fmt.Errorf("failed to write parts to output: %w", err)
		}
	case flags.OutputFormatJson:
		jb, err := json.Marshal(out)
		if err != nil {
			return fmt.Errorf("failed to serialize output to json: %w", err)
		}

		if _, err := fmt.Fprintln(cmd.OutOrStdout(), string(jb)); err != nil {
			return fmt.Errorf("failed to write parts to output: %w", err)
		}
	default:
		return fmt.Errorf("failed to format in requested format, %s is not supported", persistentFlags.OutputFormat)
	}

	return nil
}
