/*
Copyright © 2023 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"crypto/rand"
	"crypto/x509/pkix"
	"encoding/pem"
	"os"

	"github.com/spf13/cobra"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/x509"
)

// csrCmd represents the csr command
var csrCmd = &cobra.Command{
	Use:   "csr",
	Short: "generate csr",
	Long:  `generate csr`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return generateCsr()
	},
}

func generateCsr() error {
	privateKey, err := sm2.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	// 将私钥保存到文件
	privateKeyBytes, err := x509.MarshalSm2PrivateKey(privateKey, nil)
	if err != nil {
		return err
	}
	privateKeyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	pubKey, _ := privateKey.Public().(*sm2.PublicKey)
	publicKeyPem, err := x509.WritePublicKeyToPem(pubKey)
	if err != nil {
		return err
	}
	err = os.WriteFile("node.pub", publicKeyPem, 0o755)
	if err != nil {
		return err
	}

	err = os.WriteFile("node.key", privateKeyPem, 0o755)
	if err != nil {
		return err
	}

	// 创建证书签名请求模板
	template := x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName:         "node1",
			Organization:       []string{"hyperchain"},
			OrganizationalUnit: []string{"ecert"},
		},
		PublicKeyAlgorithm: x509.PublicKeyAlgorithm(x509.SM2WithSM3),
	}

	// 生成证书签名请求
	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &template, privateKey)
	if err != nil {
		return err
	}

	// 将证书签名请求保存到文件
	csrPem := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrBytes,
	})

	err = os.WriteFile("node.csr", csrPem, 0o755)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(csrCmd)
}
