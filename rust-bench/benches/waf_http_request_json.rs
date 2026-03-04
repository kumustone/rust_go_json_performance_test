use criterion::{criterion_group, criterion_main, Throughput};
use waf_http_request_serde_bench::{build_bench_request, WafHttpRequest, BENCH_CASES};

fn bench_waf_http_request_json(c: &mut criterion::Criterion) {
    let mut group = c.benchmark_group("waf_http_request_json");

    for tc in BENCH_CASES {
        let req = build_bench_request(tc);
        let payload = serde_json::to_vec(&req).expect("serialize fixture");

        group.throughput(Throughput::Bytes(payload.len() as u64));

        group.bench_function(format!("{}/encode/serde_json", tc.name), |b| {
            b.iter(|| {
                let encoded = serde_json::to_vec(&req).expect("serde_json::to_vec failed");
                criterion::black_box(encoded);
            })
        });

        group.bench_function(format!("{}/decode/serde_json", tc.name), |b| {
            b.iter(|| {
                let decoded: WafHttpRequest =
                    serde_json::from_slice(&payload).expect("serde_json::from_slice failed");
                criterion::black_box(decoded);
            })
        });
    }

    group.finish();
}

criterion_group!(benches, bench_waf_http_request_json);
criterion_main!(benches);
