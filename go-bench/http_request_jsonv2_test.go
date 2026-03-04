
// JSON v2 benchmark - 需要手动安装 github.com/go-json-experiment/json
// 由于该包仍在实验阶段，默认不包含在基准测试中
// 
// 安装方式：
// go get github.com/go-json-experiment/json@latest
//
// 运行方式：
// go test -tags jsonv2 -bench BenchmarkHttpRequestJSONv2 -benchmem

package benchmark

import (
	"testing"

	jsonv2 "github.com/go-json-experiment/json"
)

func BenchmarkHttpRequestJSONv2(b *testing.B) {
	for _, tc := range httpRequestBenchCases {
		tc := tc
		req := buildBenchRequest(tc)

		// 预生成 payload
		v2Payload, err := jsonv2.Marshal(req)
		if err != nil {
			b.Fatalf("prepare jsonv2 payload failed: %v", err)
		}

		b.Run(tc.name, func(b *testing.B) {
			b.Run("encode/jsonv2", func(b *testing.B) {
				b.ReportAllocs()
				b.SetBytes(int64(len(v2Payload)))
				for i := 0; i < b.N; i++ {
					payload, err := jsonv2.Marshal(req)
					if err != nil {
						b.Fatalf("jsonv2.Marshal failed: %v", err)
					}
					if len(payload) == 0 {
						b.Fatal("empty payload")
					}
				}
			})

			b.Run("decode/jsonv2", func(b *testing.B) {
				b.ReportAllocs()
				b.SetBytes(int64(len(v2Payload)))
				for i := 0; i < b.N; i++ {
					var out HttpRequest
					if err := jsonv2.Unmarshal(v2Payload, &out); err != nil {
						b.Fatalf("jsonv2.Unmarshal failed: %v", err)
					}
					if len(out.Body) == 0 {
						b.Fatal("empty body")
					}
				}
			})
		})
	}
}
