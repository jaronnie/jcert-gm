# jcert-gm
基于国密 sm2 算法的自签证书工具, 并保存为 pkcs7 格式.

## Usage

```shell
jcert-gm completion zsh > "${fpath[1]}/_jcert-gm" # 命令自动补全
jcert-gm init                                     # 初始化机构的根 CA, 需要保存好
jcert-gm csr                                      # 生成 privateKey 和 csr
jcert-gm cert                                     # 根据 csr 生成 cert
```

## example

```shell
./generate.sh
```

将会生成 4 个节点证书, 1 个 sdk 证书, 1 个 tls 证书

## todo

- [ ] jcert-gm parse
- [ ] jcert-gm scope
- [ ] jcert-gm print
- [ ] jcert-gm match

## 鸣谢

* [github.com/tjfoc/gmsm](github.com/tjfoc/gmsm)

