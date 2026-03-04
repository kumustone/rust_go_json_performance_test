package benchmark

import (
	"encoding/json"
	"fmt"
	"testing"
)

type httpRequestBenchCase struct {
	name            string
	bodySize        int
	headerKeyCount  int
	headerValueEach int
}

var httpRequestBenchCases = []httpRequestBenchCase{
	{name: "S", bodySize: 512, headerKeyCount: 6, headerValueEach: 1},
	{name: "M", bodySize: 8 * 1024, headerKeyCount: 12, headerValueEach: 2},
	{name: "L", bodySize: 64 * 1024, headerKeyCount: 20, headerValueEach: 3},
}

type stdHttpRequest struct {
	Mark          string
	Method        string
	Scheme        string
	Url           string
	Proto         string
	Host          string
	RemoteAddr    string
	ContentLength uint64
	Header        map[string][]string
	Body          []byte
}

func toStdHttpRequest(req HttpRequest) stdHttpRequest {
	return stdHttpRequest{
		Mark:          req.Mark,
		Method:        req.Method,
		Scheme:        req.Scheme,
		Url:           req.Url,
		Proto:         req.Proto,
		Host:          req.Host,
		RemoteAddr:    req.RemoteAddr,
		ContentLength: req.ContentLength,
		Header:        req.Header,
		Body:          req.Body,
	}
}

func buildBenchRequest(tc httpRequestBenchCase) HttpRequest {
	body := makeBenchBody(tc.bodySize)

	return HttpRequest{
		Mark:          "www.example.com",
		Method:        "POST",
		Scheme:        "https",
		Url:           "/api/v1/login?from=bench",
		Proto:         "HTTP/1.1",
		Host:          "www.example.com",
		RemoteAddr:    "10.1.2.3",
		ContentLength: uint64(len(body)),
		Header:        makeBenchHeader(tc.headerKeyCount, tc.headerValueEach),
		Body:          body,
	}
}

func makeBenchHeader(keyCount, valueEach int) map[string][]string {
	header := make(map[string][]string, keyCount)
	for i := 0; i < keyCount; i++ {
		key := fmt.Sprintf("X-Bench-%02d", i)
		values := make([]string, 0, valueEach)
		for j := 0; j < valueEach; j++ {
			values = append(values, fmt.Sprintf("v-%02d-%02d", i, j))
		}
		header[key] = values
	}
	return header
}

func makeBenchBody(size int) []byte {
	body := make([]byte, size)
	const alphabet = "abcdefghijklmnopqrstuvwxyz0123456789"
	for i := range body {
		body[i] = alphabet[i%len(alphabet)]
	}
	return body
}

func BenchmarkHttpRequestJSON(b *testing.B) {
	for _, tc := range httpRequestBenchCases {
		tc := tc
		req := buildBenchRequest(tc)
		stdReq := toStdHttpRequest(req)

		easyPayload, err := req.MarshalJSON()
		if err != nil {
			b.Fatalf("prepare easyjson payload failed: %v", err)
		}
		stdPayload, err := json.Marshal(stdReq)
		if err != nil {
			b.Fatalf("prepare stdlib payload failed: %v", err)
		}

		b.Run(tc.name, func(b *testing.B) {
			b.Run("encode/easyjson", func(b *testing.B) {
				b.ReportAllocs()
				b.SetBytes(int64(len(easyPayload)))
				for i := 0; i < b.N; i++ {
					payload, err := req.MarshalJSON()
					if err != nil {
						b.Fatalf("MarshalJSON failed: %v", err)
					}
					if len(payload) == 0 {
						b.Fatal("empty payload")
					}
				}
			})

			b.Run("encode/stdlib", func(b *testing.B) {
				b.ReportAllocs()
				b.SetBytes(int64(len(stdPayload)))
				for i := 0; i < b.N; i++ {
					payload, err := json.Marshal(stdReq)
					if err != nil {
						b.Fatalf("json.Marshal failed: %v", err)
					}
					if len(payload) == 0 {
						b.Fatal("empty payload")
					}
				}
			})

			b.Run("decode/easyjson", func(b *testing.B) {
				b.ReportAllocs()
				b.SetBytes(int64(len(easyPayload)))
				for i := 0; i < b.N; i++ {
					var out HttpRequest
					if err := out.UnmarshalJSON(easyPayload); err != nil {
						b.Fatalf("UnmarshalJSON failed: %v", err)
					}
					if len(out.Body) == 0 {
						b.Fatal("empty body")
					}
				}
			})

			b.Run("decode/stdlib", func(b *testing.B) {
				b.ReportAllocs()
				b.SetBytes(int64(len(easyPayload)))
				for i := 0; i < b.N; i++ {
					var out stdHttpRequest
					if err := json.Unmarshal(easyPayload, &out); err != nil {
						b.Fatalf("json.Unmarshal failed: %v", err)
					}
					if len(out.Body) == 0 {
						b.Fatal("empty body")
					}
				}
			})
		})
	}
}
