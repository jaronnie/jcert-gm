# jcert-gm
基于国密 sm2 算法的自签证书工具, 并保存为 pkcs7 格式.

## Usage

```shell
jcert init       # 初始化机构的根 CA, 需要保存好
jcert csr        # 生成 privateKey 和 csr
jcert cert       # 根据 csr 生成 cert
jcert match      # 证书和私钥是否匹配
```

## 鸣谢

* [github.com/tjfoc/gmsm](github.com/tjfoc/gmsm)

