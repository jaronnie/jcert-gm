/*
Copyright © 2023 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"crypto/rand"
	"encoding/pem"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/x509"
)

// matchCmd represents the match command
var matchCmd = &cobra.Command{
	Use:   "match",
	Short: "check cert and privatekey is match",
	Long:  `check cert and privatekey is match`,
	Args:  cobra.ExactArgs(2),
	RunE:  match,
}

/*
	验证 csr 与 证书是否匹配:
	通常情况下，判断 CSR 和证书是否匹配的方法会比较简单，只需要判断它们对应的公钥是否相同即可。这是因为 CSR 是用来请求颁发证书的签名请求，而证书本身就是由一个已经被信任的 CA 机构签署并包含了公钥信息的文件。
	例如，在 Go 语言的标准库中，可以使用 x509.CreateCertificate 函数生成证书时，其中一个参数就是要与证书关联的 CSR，这个 CSR 的公钥会被嵌入到生成的证书中。因此，在验证证书时只需要比较 CSR 的公钥和证书的公钥是否一致即可。
	但是在某些情况下，可能需要更严格的验证，例如验证 CSR 和证书的主题信息、扩展信息和签名等是否完全匹配。

	验证证书与私钥是否匹配:
	对消息进行签名，并使用公钥验证签名以及私钥和证书是否匹配
*/

func match(cmd *cobra.Command, args []string) error {
	certF := args[0]
	tf := args[1]

	// 解码文件
	certFile, err := os.ReadFile(certF)
	if err != nil {
		return err
	}
	certBlock, _ := pem.Decode(certFile)
	if certBlock == nil {
		return errors.New("cert block is nil")
	}
	if certBlock != nil && certBlock.Type != "CERTIFICATE" {
		return errors.New("type is not CERTIFICATE")
	}
	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return err
	}

	// 解码文件
	keyPem, err := os.ReadFile(tf)
	if err != nil {
		return err
	}
	keyBlock, _ := pem.Decode(keyPem)
	if keyBlock == nil {
		return errors.New("key block is nil")
	}
	if keyBlock.Type != "EC PRIVATE KEY" {
		return errors.New("type is not EC PRIVATE KEY")
	}
	pk, err := x509.ReadPrivateKeyFromPem(keyPem, nil)
	if err != nil {
		return err
	}

	pubPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: cert.RawSubjectPublicKeyInfo})

	pub, err := x509.ReadPublicKeyFromPem(pubPEM)
	if err != nil {
		return err
	}

	// 原理是:
	// 如果签名结果 r 和 s 是由正确的私钥对消息 "sign" 生成的
	// 那么在使用相应的公钥对消息 "sign" 和签名结果 r 和 s 进行验证时，将会得到一个布尔值为 true 的结果，
	// 表示签名有效。而如果签名结果不是由正确的私钥生成，那么验证将失败，并返回一个布尔值为 false。
	// 因此，我们可以使用这个布尔值来判断签名是否有效。

	r, s, err := sm2.Sm2Sign(pk, []byte("sign"), nil, rand.Reader)
	if err != nil {
		return err
	}

	b := sm2.Sm2Verify(pub, []byte("sign"), nil, r, s)
	fmt.Println(b)

	return nil
}

func init() {
	rootCmd.AddCommand(matchCmd)
}
