# jcert-gm
基于国密 sm2 算法的自签证书工具, 并保存为 pkcs7 格式.

## Usage

```shell
jcert-gm init       # 初始化机构的根 CA, 需要保存好
jcert-gm csr        # 生成 privateKey 和 csr
jcert-gm cert       # 根据 csr 生成 cert
jcert-gm match      # 证书和私钥是否匹配
jecet-gm print      # 打印证书内容
jcert-gm parse      # 解析证书成对应结构
```

## 鸣谢

* [github.com/tjfoc/gmsm](github.com/tjfoc/gmsm)

