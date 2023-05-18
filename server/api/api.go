package api

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/emmansun/gmsm/pkcs7"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/tjfoc/gmsm/x509"
)

func Router(rg *gin.RouterGroup) {
	rg.POST("/upload", handleUpload)
	rg.GET("/download/:filename", handleDownload)
}

func handleUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		return
	}
	f, err := file.Open()
	if err != nil {
		return
	}
	fileBytes, err := io.ReadAll(f)
	if err != nil {
		return
	}
	tarfileFp := filepath.Join("data", file.Filename)
	err = os.WriteFile(tarfileFp, fileBytes, 0o755)
	if err != nil {
		return
	}
	// 解压
	tid := uuid.New().String()
	err = UnpackTarGz(tarfileFp, filepath.Join("data", tid))
	if err != nil {
		return
	}
	s, err := GetDirAllFilePathWithSuffix(filepath.Join("data", tid), ".csr")
	if err != nil {
		return
	}

	// 生成证书
	oid := uuid.New().String()
	_ = os.MkdirAll(filepath.Join("data", oid), 0o755)
	for _, v := range s {
		err = generateCert(v, filepath.Join("data", oid))
		if err != nil {
			return
		}
	}
	err = CompressFolder(filepath.Join("data", oid), filepath.Join("data", fmt.Sprintf("%s.zip", oid)))
	if err != nil {
		return
	}

	c.JSON(200, struct {
		Filename string
	}{
		Filename: fmt.Sprintf("%s.zip", oid),
	})
}

func handleDownload(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	c.File(filepath.Join("data", c.Params.ByName("filename")))
}

func generateCert(csrfp string, output string) error {
	configDir := filepath.Dir(viper.ConfigFileUsed())

	// 读取CSR文件
	csrPEM, err := os.ReadFile(csrfp)
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

	// 创建证书模板
	template := &x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               csr.Subject,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		SubjectKeyId:          []byte{1, 2, 3, 4, 6},
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageCodeSigning, x509.ExtKeyUsageEmailProtection},
		PublicKeyAlgorithm:    csr.PublicKeyAlgorithm,
		SignatureAlgorithm:    csr.SignatureAlgorithm,
		DNSNames:              csr.DNSNames,
		CRLDistributionPoints: viper.GetStringSlice("CRLDistributionPoints"),
		OCSPServer:            viper.GetStringSlice("OCSPServer"),
	}

	// 使用SM2密钥对签名证书
	derBytes, err := x509.CreateCertificate(template, ca, pub, privateKey)
	if err != nil {
		return err
	}

	// 将证书转换为PEM格式
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	generatedCert := filepath.Join(output, fmt.Sprintf("%s-%s-%s.p7b", csr.Subject.CommonName, csr.Subject.OrganizationalUnit[0], uuid.New().String()[:6]))
	p7b, err := saveCertToPkcs7(certPEM, caPEM)
	if err != nil {
		return err
	}

	err = os.WriteFile(generatedCert, p7b, 0o755)
	return err
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

// GetDirAllFilePathWithSuffix gets all the file paths in the specified directory recursively with suffix.
func GetDirAllFilePathWithSuffix(dirname string, suffix string) ([]string, error) {
	filePaths, err := GetDirAllFilePath(dirname)
	if err != nil {
		return nil, err
	}
	paths := make([]string, 0)
	for _, v := range filePaths {
		if filepath.Ext(v) == suffix {
			paths = append(paths, v)
		}
	}
	return paths, nil
}

// GetDirAllFilePath gets all the file paths in the specified directory recursively.
func GetDirAllFilePath(dirname string) ([]string, error) {
	// Remove the trailing path separator if dirname has.
	dirname = strings.TrimSuffix(dirname, string(os.PathSeparator))
	infos, err := os.ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	paths := make([]string, 0, len(infos))
	for _, info := range infos {
		path := dirname + string(os.PathSeparator) + info.Name()
		realInfo, err := os.Stat(path)
		if err != nil {
			return nil, err
		}
		if realInfo.IsDir() {
			tmp, err := GetDirAllFilePath(path)
			if err != nil {
				return nil, err
			}
			paths = append(paths, tmp...)
			continue
		}
		paths = append(paths, path)
	}
	return paths, nil
}

func UnpackTarGz(filename string, dest string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return err
		}

		path := filepath.Join(dest, hdr.Name)

		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path, hdr.FileInfo().Mode()); err != nil {
				return err
			}
		case tar.TypeReg:
			_ = os.MkdirAll(filepath.Dir(path), hdr.FileInfo().Mode())
			file, err := os.Create(path)
			if err != nil {
				return err
			}
			defer file.Close()
			if _, err := io.Copy(file, tr); err != nil {
				return err
			}
			if err := file.Chmod(hdr.FileInfo().Mode()); err != nil {
				return err
			}
		}
	}

	return nil
}

func CompressFolder(inputFolderPath string, outputZipPath string) error {
	outputFile, err := os.Create(outputZipPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	zipWriter := zip.NewWriter(outputFile)
	defer zipWriter.Close()

	// 递归地遍历输入文件夹中的所有文件和子文件夹，并将它们添加到 zip 文件中
	err = filepath.Walk(inputFolderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 获取相对路径（相对于输入文件夹）
		relPath, err := filepath.Rel(inputFolderPath, path)
		if err != nil {
			return err
		}

		// 如果是文件夹，则跳过
		if info.IsDir() {
			return nil
		}

		// 创建一个新的 zip 文件头
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = relPath

		// 将文件头写入 zip 文件
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		// 打开原始文件
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// 将原始文件内容复制到 zip 文件中
		_, err = io.Copy(writer, file)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	// 压缩完成
	fmt.Printf("Successfully compressed folder '%s' to '%s'\n", inputFolderPath, outputZipPath)
	return nil
}
