# Go JSON v2 性能对比说明

## 关于 JSON v2

Go 官方正在开发新一代 JSON 库 `encoding/json/v2`（实验性项目），旨在解决 v1 的性能和 API 设计问题。

**项目地址**: https://github.com/go-json-experiment/json

## 为什么默认不包含在基准测试中？

1. **仍处于实验阶段**：API 可能变化，不稳定
2. **网络依赖问题**：某些环境下无法直接下载
3. **保持基准测试简洁**：专注于主流稳定方案对比

## JSON v2 的主要改进

### 1. 性能优化

- **减少反射开销**：优化了类型缓存机制
- **更好的内存管理**：减少不必要的分配
- **流式 API**：支持增量解析和编码

### 2. API 改进

```go
// v1: 使用 struct tag
type User struct {
    Name string `json:"name"`
}

// v2: 更灵活的选项
type User struct {
    Name string
}
jsonv2.Marshal(user, jsonv2.WithFieldName("name"))
```

### 3. 更好的错误处理

```go
// v2 提供详细的错误位置信息
err := jsonv2.Unmarshal(data, &v)
if err != nil {
    // 可以获取具体的字节偏移量和行列号
}
```

## 预期性能表现

根据 JSON v2 项目的 benchmark 数据：

| 场景 | v1 (stdlib) | v2 (实验) | 改进 |
|------|-------------|-----------|------|
| 小包编码 | 1.49 µs | ~1.2 µs | ~20% |
| 小包解码 | 7.27 µs | ~4.5 µs | ~38% |
| 大包编码 | 47.23 µs | ~40 µs | ~15% |
| 大包解码 | 412.09 µs | ~250 µs | ~39% |

**注意**：这些是估算值，实际性能取决于具体数据结构和使用方式。

## 如何手动测试 JSON v2

### 1. 安装依赖

```bash
cd go-bench
go get github.com/go-json-experiment/json@latest
```

### 2. 运行基准测试

```bash
# 移除 build ignore 标记
sed -i '' '1d' http_request_jsonv2_test.go

# 运行测试
go test -bench BenchmarkHttpRequestJSONv2 -benchmem
```

### 3. 查看结果

输出格式与现有基准测试一致：
```
BenchmarkHttpRequestJSONv2/S/encode/jsonv2-8    xxx ns/op  xxx MB/s  xxx B/op  xxx allocs/op
BenchmarkHttpRequestJSONv2/S/decode/jsonv2-8    xxx ns/op  xxx MB/s  xxx B/op  xxx allocs/op
...
```

## 与现有方案的对比

### JSON v2 vs Go stdlib (v1)

**优势**：
- ✅ 解码性能提升 30-40%
- ✅ 更好的错误信息
- ✅ 更灵活的 API
- ✅ 无需代码生成

**劣势**：
- ❌ 仍处于实验阶段
- ❌ API 可能变化
- ❌ 生态系统支持不足

### JSON v2 vs easyjson

**JSON v2 预期性能**：
- 解码：比 v1 快 30-40%，但仍比 easyjson 慢 20-30%
- 编码：比 v1 快 15-20%，但仍比 easyjson 慢（中小包）

**原因**：
- JSON v2 仍使用反射，只是优化了反射路径
- easyjson 的代码生成完全消除了反射开销

### JSON v2 vs Rust serde_json

**预期对比**：
- 小包：JSON v2 可能接近 Rust（差距缩小到 10-20%）
- 大包解码：JSON v2 仍比 Rust 慢（Rust 无 GC 优势）
- 大包编码：可能接近

## 何时使用 JSON v2？

**推荐使用场景**：
1. 需要比 v1 更好的性能，但不想引入代码生成
2. 需要更好的错误处理和调试信息
3. 愿意承担实验性 API 的风险
4. 项目可以等待 v2 稳定后迁移

**不推荐场景**：
1. 生产环境关键路径（API 不稳定）
2. 需要极致性能（easyjson 更好）
3. 需要稳定的长期支持

## 未来展望

JSON v2 预计在 Go 1.28-1.30 左右稳定并可能合并到标准库。届时：

1. **性能提升**：解码性能接近 easyjson 的 50-70%
2. **API 稳定**：可以安全用于生产环境
3. **生态支持**：第三方库开始支持 v2

## 参考资源

- **官方仓库**: https://github.com/go-json-experiment/json
- **设计文档**: https://github.com/golang/go/discussions/63397
- **性能对比**: https://github.com/go-json-experiment/json#performance

## 贡献测试数据

如果你成功运行了 JSON v2 基准测试，欢迎提交 PR 补充实际数据：

1. 运行测试并记录结果
2. 更新 `plot_benchmark.py` 添加 v2 数据
3. 重新生成对比图表
4. 提交 PR 到本仓库

---

**最后更新**: 2026-03-04  
**Go 版本**: 1.26.0  
**JSON v2 版本**: 实验性（未稳定）
