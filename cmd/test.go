/*
Copyright © 2023 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"encoding/pem"
	"errors"
	"fmt"
	"os"

	"github.com/emmansun/gmsm/smx509"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		f := args[0]
		keyPem, err := os.ReadFile(f)
		if err != nil {
			return err
		}
		keyBlock, _ := pem.Decode(keyPem)
		if keyBlock == nil {
			return errors.New("key block is nil")
		}
		pk, err := smx509.ParseTypedECPrivateKey(keyBlock.Bytes)
		if err != nil {
			return err
		}
		fmt.Println(pk)

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