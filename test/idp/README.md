# IdP 端到端测试

RAM 身份提供商（IdP）相关控制器、库表与服务层的端到端测试。

- 测试说明：[docs/test/idp-e2e-test.md](../../docs/test/idp-e2e-test.md)
- 测试结果：[docs/test/idp-e2e-test-result.md](../../docs/test/idp-e2e-test-result.md)

快速运行：

```bash
export GOROOT=/Users/mac/develop/go
export PATH=$PATH:$GOROOT/bin
export GOCACHE=/tmp/dsh-go-build
export MOCKEY_CHECK_GCFLAGS=false
go test ./test/idp -run TestIDP -v -count=1
```

测试连接 `data/config/app.toml` 配置的 PostgreSQL，会创建带 `E2E-` 前缀的数据并在结束时清理。
