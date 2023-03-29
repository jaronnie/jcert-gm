/*
Copyright © 2023 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"bytes"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/emmansun/gmsm/pkcs7"
	"github.com/spf13/cobra"
)

// transCmd represents the trans command
var transCmd = &cobra.Command{
	Use:   "trans",
	Short: "trans pkcs7 to pem",
	Long:  `trans pkcs7 to pem`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		f := args[0]

		b, err := os.ReadFile(f)
		cobra.CheckErr(err)

		p, err := base64.RawStdEncoding.DecodeString(string(b))
		cobra.CheckErr(err)

		// 解码 PKCS#7 数据
		certs, err := pkcs7.Parse(p)
		if err != nil {
			fmt.Println("Error parsing PKCS#7 data:", err)
			os.Exit(1)
		}

		// 将解码后的证书转换为 PEM 格式
		buffer := &bytes.Buffer{}
		for _, cert := range certs.Certificates {
			block := &pem.Block{
				Type:  "CERTIFICATE",
				Bytes: cert.Raw,
			}
			pemData := pem.EncodeToMemory(block)
			buffer.Write(pemData)
		}
		fmt.Printf("%s\n", buffer.Bytes())
	},
}

func init() {
	rootCmd.AddCommand(transCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// transCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// transCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
