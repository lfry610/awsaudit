package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type identity struct {
	Account string `json:"account"`
	ARN     string `json:"arn"`
}

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Shows current AWS account number and ARN",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		cfg, err := loadAWSConfig(ctx)
		if err != nil {
			return fmt.Errorf("Loading AWS Config: %w", err)
		}
		client := sts.NewFromConfig(cfg)
		out, err := client.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
		if err != nil {
			return fmt.Errorf("Calling STS: %w", err)
		}

		switch viper.GetString(keyOutput) {
		case "json":
			id := identity{
				Account: aws.ToString(out.Account),
				ARN:     aws.ToString(out.Arn),
			}
			enc := json.NewEncoder(cmd.OutOrStdout())
			enc.SetIndent("", "  ")
			if err = enc.Encode(id); err != nil {
				return fmt.Errorf("Encoding JSON: %w", err)
			}
		case "table":
			fmt.Fprintf(cmd.OutOrStdout(), "Account:  %s\n", aws.ToString(out.Account))
			fmt.Fprintf(cmd.OutOrStdout(), "ARN:  %s\n", aws.ToString(out.Arn))
		default:
			return fmt.Errorf("Unknown output format %q (want table or json)", viper.GetString(keyOutput))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(whoamiCmd)
}
