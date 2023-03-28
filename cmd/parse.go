/*
Copyright © 2023 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"encoding/pem"
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
	Args:  cobra.ExactArgs(1),
	RunE:  parse,
}

func parse(cmd *cobra.Command, args []string) error {
	// 读取文件
	file, err := os.ReadFile(args[0])
	if err != nil {
		return err
	}

	// 解码文件
	var t string
	block, _ := pem.Decode(file)
	if block != nil && block.Type == "CERTIFICATE REQUEST" {
		t = "csr"
	} else if block.Type == "CERTIFICATE" {
		t = "cert"
	}

	if t == "csr" {
		parse, err := x509.ParseCertificateRequest(block.Bytes)
		if err != nil {
			return err
		}
		fmt.Printf("common name: %s\n", parse.Subject.CommonName)
		fmt.Printf("Organization: %s\n", strings.Join(parse.Subject.Organization, ","))
	}

	if t == "cert" {
		parse, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return err
		}
		fmt.Printf("common name: %s\n", parse.Subject.CommonName)
		fmt.Printf("Organization: %s\n", strings.Join(parse.Subject.Organization, ","))
	}

	return nil
}

func init() {
	rootCmd.AddCommand(parseCmd)
}
