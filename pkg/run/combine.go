package run

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/shamir"
	keysreader "github.com/kubetrail/bip32/pkg/keys"
	"github.com/kubetrail/bip39/pkg/prompts"
	"github.com/kubetrail/mkchop/pkg/flags"
	"github.com/mr-tron/base58"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func Combine(cmd *cobra.Command, args []string) error {
	persistentFlags := getPersistentFlags(cmd)

	_ = viper.BindPFlag(flags.NumThreshold, cmd.Flags().Lookup(flags.NumThreshold))
	_ = viper.BindPFlag(flags.Key, cmd.Flag(flags.Key))
	numThreshold := viper.GetInt(flags.NumThreshold)
	keys := viper.GetStringSlice(flags.Key)

	if numThreshold < 2 {
		return fmt.Errorf("num-threshold has to be at least 2")
	}

	keys = append(keys, args...)
	if len(keys) > numThreshold {
		keys = keys[:numThreshold]
	}

	prompt, err := prompts.Status()
	if err != nil {
		return fmt.Errorf("failed to get prompt status: %w", err)
	}

	partsMap := make(map[string]struct{})
	parts := make([][]byte, numThreshold)
	if len(keys) < numThreshold {
		for i := len(keys); i < numThreshold; i++ {
			if prompt {
				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Enter part %d of %d: ", i+1, numThreshold); err != nil {
					return fmt.Errorf("failed to write to output: %w", err)
				}
			}

			key, err := keysreader.Read(cmd.InOrStdin())
			if err != nil {
				return fmt.Errorf("failed to read key part from input: %w", err)
			}

			keys = append(keys, key)
		}
	}

	for i, key := range keys {
		partsMap[key] = struct{}{}
		parts[i], err = base58.Decode(key)
		if err != nil {
			return fmt.Errorf("failed to base58 decode part %d: %w", i, err)
		}
	}

	if len(partsMap) != numThreshold {
		return fmt.Errorf("please input at least %d unique parts", numThreshold)
	}

	key, err := shamir.Combine(parts)
	if err != nil {
		return fmt.Errorf("failed to combine parts: %w", err)
	}

	type output struct {
		Key string `json:"key,omitempty" yaml:"key,omitempty"`
	}

	out := &output{
		Key: string(key),
	}

	switch strings.ToLower(persistentFlags.OutputFormat) {
	case flags.OutputFormatNative:
		if _, err := fmt.Fprintln(cmd.OutOrStdout(), out.Key); err != nil {
			return fmt.Errorf("failed to write part to output: %w", err)
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
