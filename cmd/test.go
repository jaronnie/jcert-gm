/*
Copyright © 2023 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"encoding/pem"
	"os"

	"github.com/tjfoc/gmsm/x509"

	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "test",
	Long:  `test`,
	RunE: func(cmd *cobra.Command, args []string) error {
		f := args[0]
		p, err := os.ReadFile(f)
		if err != nil {
			return err
		}

		var ca *x509.Certificate
		var cert *x509.Certificate

		for {
			keyBlock, rest := pem.Decode(p)
			if keyBlock == nil {
				break
			}
			c, err := x509.ParseCertificate(keyBlock.Bytes)
			if err != nil {
				return err
			}
			if c.IsCA {
				ca = c
			} else {
				cert = c
			}

			p = rest
		}

		err = ca.CheckSignatureFrom(cert)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
