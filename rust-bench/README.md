# rust-serde benchmark

Benchmark `WafHttpRequest` JSON encode/decode with `serde_json`.

## Run

```bash
cargo bench --manifest-path bench/rust-serde/Cargo.toml
```

Or from repo root:

```bash
make bench-rust
```

## Notes

- Field names follow Go JSON keys (`Mark`, `Method`, ...).
- `Body` uses base64 encode/decode to align with Go JSON behavior for `[]byte`.
