/*
Copyright © 2023 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"

	"github.com/emmansun/gmsm/pkcs7"
	"github.com/tjfoc/gmsm/x509"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Csr    string
	Output string
)

// certCmd represents the cert command
var certCmd = &cobra.Command{
	Use:   "cert",
	Short: "generate cert by csr",
	Long:  `generate cert by csr`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return generateCert()
	},
}

func generateCert() error {
	configDir := filepath.Dir(viper.ConfigFileUsed())

	// 读取CSR文件
	csrPEM, err := os.ReadFile(Csr)
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

	pubPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: csr.RawSubjectPublicKeyInfo})

	pub, err := x509.ReadPublicKeyFromPem(pubPEM)
	if err != nil {
		return err
	}

	// 申请序列号
	// 随机生成一个
	serialNumber, err := rand.Int(rand.Reader, big.NewInt(1<<63-1))
	if err != nil {
		return err
	}

	// 创建证书模板
	template := &x509.Certificate{
		SerialNumber:       serialNumber,
		Subject:            csr.Subject,
		NotBefore:          time.Now(),
		NotAfter:           time.Now().Add(365 * 24 * time.Hour),
		SubjectKeyId:       []byte{1, 2, 3, 4, 6},
		KeyUsage:           x509.KeyUsageDigitalSignature,
		ExtKeyUsage:        []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageCodeSigning, x509.ExtKeyUsageEmailProtection},
		PublicKeyAlgorithm: csr.PublicKeyAlgorithm,
		SignatureAlgorithm: csr.SignatureAlgorithm,
		DNSNames:           csr.DNSNames,
	}

	// 读取机构 ca 文件
	caPEM, err := os.ReadFile(filepath.Join(configDir, "ca.cert"))
	if err != nil {
		return err
	}

	// 解码 ca 文件
	caBlock, _ := pem.Decode(caPEM)
	if err != nil {
		return err
	}

	ca, err := x509.ParseCertificate(caBlock.Bytes)
	if err != nil {
		return err
	}

	// 读取机构 ca 私钥
	caKey, err := os.ReadFile(filepath.Join(configDir, "ca.key"))
	if err != nil {
		return err
	}

	privateKey, err := x509.ReadPrivateKeyFromPem(caKey, nil)
	if err != nil {
		return err
	}

	// 使用SM2密钥对签名证书
	derBytes, err := x509.CreateCertificate(template, ca, pub, privateKey)
	if err != nil {
		return err
	}

	// 将证书转换为PEM格式
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	if Output == "pem" {
		generatedCert := filepath.Join(Path, fmt.Sprintf("%s.cert", csr.Subject.CommonName))
		pem := savaCertToPem(certPEM, caPEM)
		_ = os.WriteFile(generatedCert, pem, 0o755)
		return nil
	}

	if Output == "pkcs7" {
		generatedCert := filepath.Join(Path, fmt.Sprintf("%s.p7b", csr.Subject.CommonName))
		p7b, err := saveCertToPkcs7(certPEM, caPEM)
		if err != nil {
			return err
		}

		_ = os.WriteFile(generatedCert, p7b, 0o755)
		return nil
	}
	return errors.Errorf("not suuport output %s", Output)
}

func saveCertToPkcs7(cert []byte, ca []byte) ([]byte, error) {
	b := savaCertToPem(cert, ca)

	buffer := &bytes.Buffer{}
	for {
		block, rest := pem.Decode(b)
		if block == nil {
			break
		}
		buffer.Write(block.Bytes)
		b = rest
	}

	p7b, err := pkcs7.DegenerateCertificate(buffer.Bytes())
	if err != nil {
		return nil, err
	}

	s := base64.StdEncoding.EncodeToString(p7b)

	return []byte(s), nil
}

func savaCertToPem(cert []byte, ca []byte) []byte {
	buffer := &bytes.Buffer{}

	buffer.Write(ca)
	buffer.Write([]byte("\n"))
	buffer.Write(cert)

	return buffer.Bytes()
}

func init() {
	rootCmd.AddCommand(certCmd)

	certCmd.Flags().StringVarP(&Csr, "csr", "", "", "set csr file path")
	certCmd.Flags().StringVarP(&Output, "output", "o", "pem", "set cert block type")

	_ = certCmd.MarkFlagRequired("csr")
}
