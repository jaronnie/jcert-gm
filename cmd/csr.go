/*
Copyright © 2023 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"crypto/rand"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"

	ssm2 "github.com/emmansun/gmsm/sm2"
	"github.com/emmansun/gmsm/smx509"
	"github.com/spf13/cobra"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/x509"
)

/*

	csr 即证书签名请求, 可以通过 csr 向 ca 机构申请证书, 在本工具中可以通过执行 init 命令, 模拟了 ca 机构创建出机构的根 ca 和 私钥
	可通过 csr 命令生成 csr 文件后, 再执行 cert 命令将 csr 作为参数传进去, 即可申请到证书.

	通过本命令将生成三个文件:
	1. 私钥, 首先生成私钥, 只支持使用 sm2 国密算法
	2. 公钥, 从生成的私钥中取出公钥, 保存在文件中
	3. 根据私钥生成 csr 证书签名文件, 可选择签名算法, 默认只支持 sm2-sha256

*/

var (
	CN   string
	O    []string
	OU   []string
	Addr []string

	EC bool
)

// csrCmd represents the csr command
var csrCmd = &cobra.Command{
	Use:   "csr",
	Short: "generate csr",
	Long:  `generate csr`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if CN == "" {
			cobra.CheckErr("cn is empty")
		}
		return generateCsr()
	},
}

func generateCsr() error {
	var (
		generatedKey = filepath.Join(Path, fmt.Sprintf("%s.key", CN))
		generatedPub = filepath.Join(Path, fmt.Sprintf("%s.pub", CN))
		generatedCsr = filepath.Join(Path, fmt.Sprintf("%s.csr", CN))
	)

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
		Type:  "PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	pubKey, _ := privateKey.Public().(*sm2.PublicKey)
	publicKeyPem, err := x509.WritePublicKeyToPem(pubKey)
	if err != nil {
		return err
	}
	err = os.WriteFile(generatedPub, publicKeyPem, 0o755)
	if err != nil {
		return err
	}

	if EC {
		ecPrivateKeyPem, err := transPkcs8PrivateKeyToEcPrivateKey(privateKeyBytes)
		if err != nil {
			return err
		}
		err = os.WriteFile(generatedKey, ecPrivateKeyPem, 0o755)
		if err != nil {
			return err
		}
	} else {
		err = os.WriteFile(generatedKey, privateKeyPem, 0o755)
		if err != nil {
			return err
		}
	}

	// 创建证书签名请求模板
	template := x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName:         CN,
			Organization:       O,
			OrganizationalUnit: OU,
			Province:           []string{"浙江省"},
			Locality:           []string{"杭州市"},
			Country:            []string{"CN"},
		},
		SignatureAlgorithm: x509.SM2WithSM3,
		PublicKeyAlgorithm: x509.PublicKeyAlgorithm(x509.SM2WithSM3),
		DNSNames:           Addr,
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

	err = os.WriteFile(generatedCsr, csrPem, 0o755)
	if err != nil {
		return err
	}
	return nil
}

func transPkcs8PrivateKeyToEcPrivateKey(der []byte) ([]byte, error) {
	// 将 pkcs8 格式的私钥转化成 ecPrivateKey
	key, err := smx509.ParsePKCS8PrivateKey(der)
	if err != nil {
		return nil, err
	}
	b, err := smx509.MarshalSM2PrivateKey(key.(*ssm2.PrivateKey))
	if err != nil {
		return nil, err
	}

	// 将证书签名请求保存到文件
	pem := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: b,
	})

	return pem, nil
}

func init() {
	rootCmd.AddCommand(csrCmd)

	csrCmd.Flags().StringVarP(&CN, "CN", "", "", "set CommonName")
	csrCmd.Flags().StringSliceVarP(&O, "O", "", nil, "set Organization")
	csrCmd.Flags().StringSliceVarP(&OU, "OU", "", nil, "set OrganizationUnit")
	csrCmd.Flags().StringSliceVarP(&Addr, "addr", "", nil, "set dns addr")

	csrCmd.Flags().StringVarP(&Path, "path", "p", "", "save path")

	csrCmd.Flags().BoolVarP(&EC, "ec", "", false, "trans pkcs8 private key to ec private key")

	_ = csrCmd.MarkFlagRequired("CN")
}
