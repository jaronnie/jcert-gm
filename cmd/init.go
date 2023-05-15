/*
Copyright © 2023 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"crypto/rand"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/x509"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the root ca of the mechanism",
	Long:  `Initialize the root ca of the mechanism`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return generateAuthorityRootCA()
	},
}

func generateAuthorityRootCA() error {
	configDir := filepath.Dir(viper.ConfigFileUsed())

	// 创建 CA私钥
	caPrivKey, err := sm2.GenerateKey(rand.Reader)
	if err != nil {
		return err
	}

	// 将私钥保存到文件
	caPrivateKeyBytes, err := x509.MarshalSm2PrivateKey(caPrivKey, nil)
	if err != nil {
		return err
	}
	caPrivateKeyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: caPrivateKeyBytes,
	})
	err = os.WriteFile(filepath.Join(configDir, "ca.key"), caPrivateKeyPem, 0o755)
	if err != nil {
		return err
	}

	// 创建 CA 证书模板
	caTemplate := x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			CommonName:   "BLOCFACE HYPERCHAIN SM2 OCA1",
			Organization: []string{"Blocface Hyperchain Self Authority"},
			Country:      []string{"CN"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(100, 0, 0), // 有效期为 100 年
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	// 创建自签的 CA 证书
	caDerBytes, err := x509.CreateCertificate(&caTemplate, &caTemplate, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return err
	}
	b := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: caDerBytes})

	err = os.WriteFile(filepath.Join(configDir, "ca.cert"), b, 0o755)
	if err != nil {
		return err
	}

	// create crl
	crlBytes, err := caTemplate.CreateCRL(rand.Reader, caPrivKey, nil, time.Now(), time.Now().Add(time.Duration(time.Now().Year())*100))
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(configDir, "crl.crl"), crlBytes, 0o755)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(initCmd)
}
