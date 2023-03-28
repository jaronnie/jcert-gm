/*
Copyright © 2023 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/tjfoc/gmsm/x509"

	"github.com/spf13/cobra"
)

// parseCmd represents the parse command
var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "parse csr or certs",
	Long:  `parse cer or certs`,
	Args:  cobra.ExactArgs(2),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"csr", "cert"}, cobra.ShellCompDirectiveDefault
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return parse(cmd, args)
	},
}

func parse(cmd *cobra.Command, args []string) error {
	// 分为 csr 和 cert
	t := args[0]
	f := args[1]

	if t == "csr" {
		// 读取CSR文件
		csrPEM, err := os.ReadFile(f)
		if err != nil {
			return err
		}

		// 解码CSR文件
		csrBlock, _ := pem.Decode(csrPEM)
		if csrBlock == nil || csrBlock.Type != "CERTIFICATE REQUEST" {
			return errors.New("type is not CERTIFICATE REQUEST")
		}
		csr, err := x509.ParseCertificateRequest(csrBlock.Bytes)
		if err != nil {
			return err
		}
		fmt.Printf("common name: %s\n", csr.Subject.CommonName)
		fmt.Printf("Organization: %s\n", strings.Join(csr.Subject.Organization, ","))
	}

	if t == "cert" {
		// 读取机构 ca 文件
		certPEM, err := os.ReadFile(f)
		if err != nil {
			return err
		}

		// 解码 ca 文件
		certBlock, _ := pem.Decode(certPEM)
		if err != nil {
			return err
		}
		cert, err := x509.ParseCertificate(certBlock.Bytes)
		if err != nil {
			return err
		}
		fmt.Printf("common name: %s\n", cert.Subject.CommonName)
		fmt.Printf("Organization: %s\n", strings.Join(cert.Subject.Organization, ","))
	}

	return nil
}

func init() {
	rootCmd.AddCommand(parseCmd)
}
