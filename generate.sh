jcert-gm csr --CN node1 --O hyperchain --OU ecert -p testdata/node1
jcert-gm cert --csr testdata/node1/node1.csr -p testdata/node1
jcert-gm csr --CN node2 --O hyperchain --OU ecert -p testdata/node2
jcert-gm cert --csr testdata/node2/node2.csr -p testdata/node2
jcert-gm csr --CN node3 --O hyperchain --OU ecert -p testdata/node3
jcert-gm cert --csr testdata/node3/node3.csr -p testdata/node3
jcert-gm csr --CN node4 --O hyperchain --OU ecert -p testdata/node4
jcert-gm cert --csr testdata/node4/node4.csr -p testdata/node4
jcert-gm csr --CN sdk.cn --O hyperchain --OU sdkcert -p testdata/sdk
jcert-gm cert --csr testdata/sdk/sdk.cn.csr -p testdata/sdk
jcert-gm csr --CN hyperchain.cn --O hyperchain --addr hyperchain.cn -p testdata/tls
jcert-gm cert --csr testdata/tls/hyperchain.cn.csr -p testdata/tls