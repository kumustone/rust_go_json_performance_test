.PHONY: all bench-go bench-rust plot clean setup

all: bench-go bench-rust plot

setup:
	@echo "Setting up Go benchmark..."
	@cd go-bench && go mod download
	@cd go-bench && go run github.com/mailru/easyjson/easyjson@latest -all http_request.go
	@echo "Go setup complete."
	@echo ""
	@echo "Setting up Rust benchmark..."
	@cd rust-bench && cargo build --release
	@echo "Rust setup complete."

bench-go:
	@echo "Running Go benchmarks..."
	@cd go-bench && go test -run '^$$' -bench BenchmarkHttpRequestJSON -benchmem

bench-rust:
	@echo "Running Rust benchmarks..."
	@cd rust-bench && cargo bench

plot:
	@echo "Generating comparison plots..."
	@python3 plot_benchmark.py

clean:
	@rm -f benchmark_comparison.png benchmark_ratio.png
	@cd go-bench && rm -f http_request_easyjson.go go.sum
	@cd rust-bench && cargo clean
