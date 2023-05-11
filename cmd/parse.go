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

	"github.com/fatih/color"

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

func parse(_ *cobra.Command, args []string) error {
	// 读取文件
	file, err := os.ReadFile(args[0])
	if err != nil {
		return err
	}

	count := 0

	for {
		// 解码文件
		block, rest := pem.Decode(file)
		if block == nil {
			break
		}
		file = rest

		if count >= 1 {
			fmt.Printf("\n===================================\n")
		}

		if block.Type == "CERTIFICATE REQUEST" {
			parse, err := x509.ParseCertificateRequest(block.Bytes)
			if err != nil {
				return err
			}
			fmt.Println(color.BlueString("CERTIFICATE REQUEST\n"))
			fmt.Println(color.CyanString("Subject:\n"))
			fmt.Printf("common name: %s\n", parse.Subject.CommonName)
			fmt.Printf("Organization: %s\n", strings.Join(parse.Subject.Organization, ","))
			fmt.Printf("Organization Unit: %s\n", strings.Join(parse.Subject.OrganizationalUnit, ","))
		} else if block.Type == "CERTIFICATE" {
			parse, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return err
			}
			fmt.Println(color.BlueString("CERTIFICATE\n"))
			fmt.Println(color.CyanString("Issuer:\n"))
			fmt.Printf("common name: %s\n", parse.Issuer.CommonName)
			fmt.Printf("Organization: %s\n", strings.Join(parse.Issuer.Organization, ","))
			fmt.Printf("Organization Unit: %s\n", strings.Join(parse.Issuer.OrganizationalUnit, ","))

			fmt.Println(color.CyanString("\nSubject:\n"))
			fmt.Printf("common name: %s\n", parse.Subject.CommonName)
			fmt.Printf("Organization: %s\n", strings.Join(parse.Subject.Organization, ","))
			fmt.Printf("Organization Unit: %s\n", strings.Join(parse.Subject.OrganizationalUnit, ","))
		}
		count++
	}
	return nil
}

func init() {
	rootCmd.AddCommand(parseCmd)
}
