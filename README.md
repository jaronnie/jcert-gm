# jcert-gm

基于国密 sm2 算法的自签证书工具, 并保存为 pkcs7 格式.

## Usage

```shell
jcert-gm completion zsh > "${fpath[1]}/_jcert-gm" # 命令自动补全
jcert-gm init                                     # 初始化机构的根 CA, 需要保存好
jcert-gm csr                                      # 生成 privateKey 和 csr
jcert-gm cert                                     # 根据 csr 生成 cert
jcert-gm match cert key                           # 检查私钥和证书是否匹配
```

## todo

- [ ] jcert-gm scope

## 鸣谢

- [github.com/tjfoc/gmsm](https://github.com/tjfoc/gmsm)
- [github.com/emmansun/gmsm](https://github.com/emmansun/gmsm)

