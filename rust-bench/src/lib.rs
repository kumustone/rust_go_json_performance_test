use base64::engine::general_purpose::STANDARD;
use base64::Engine;
use rustc_hash::FxHashMap;
use serde::de::Error as DeError;
use serde::{Deserialize, Deserializer, Serialize, Serializer};

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct WafHttpRequest {
    #[serde(rename = "Mark")]
    pub mark: String,
    #[serde(rename = "Method")]
    pub method: String,
    #[serde(rename = "Scheme")]
    pub scheme: String,
    #[serde(rename = "Url")]
    pub url: String,
    #[serde(rename = "Proto")]
    pub proto: String,
    #[serde(rename = "Host")]
    pub host: String,
    #[serde(rename = "RemoteAddr")]
    pub remote_addr: String,
    #[serde(rename = "ContentLength")]
    pub content_length: u64,
    #[serde(rename = "Header")]
    pub header: FxHashMap<String, Vec<String>>,
    #[serde(
        rename = "Body",
        serialize_with = "serialize_body_base64",
        deserialize_with = "deserialize_body_base64"
    )]
    pub body: Vec<u8>,
}

#[derive(Clone, Copy, Debug)]
pub struct BenchCase {
    pub name: &'static str,
    pub body_size: usize,
    pub header_key_count: usize,
    pub header_value_each: usize,
}

pub const BENCH_CASES: [BenchCase; 3] = [
    BenchCase {
        name: "S",
        body_size: 512,
        header_key_count: 6,
        header_value_each: 1,
    },
    BenchCase {
        name: "M",
        body_size: 8 * 1024,
        header_key_count: 12,
        header_value_each: 2,
    },
    BenchCase {
        name: "L",
        body_size: 64 * 1024,
        header_key_count: 20,
        header_value_each: 3,
    },
];

pub fn build_bench_request(tc: BenchCase) -> WafHttpRequest {
    let body = make_bench_body(tc.body_size);

    WafHttpRequest {
        mark: "www.example.com".to_owned(),
        method: "POST".to_owned(),
        scheme: "https".to_owned(),
        url: "/api/v1/login?from=bench".to_owned(),
        proto: "HTTP/1.1".to_owned(),
        host: "www.example.com".to_owned(),
        remote_addr: "10.1.2.3".to_owned(),
        content_length: body.len() as u64,
        header: make_bench_header(tc.header_key_count, tc.header_value_each),
        body,
    }
}

fn make_bench_header(key_count: usize, value_each: usize) -> FxHashMap<String, Vec<String>> {
    let mut header = FxHashMap::with_capacity_and_hasher(key_count, Default::default());
    for i in 0..key_count {
        let key = format!("X-Bench-{i:02}");
        let mut values = Vec::with_capacity(value_each);
        for j in 0..value_each {
            values.push(format!("v-{i:02}-{j:02}"));
        }
        header.insert(key, values);
    }
    header
}

fn make_bench_body(size: usize) -> Vec<u8> {
    let alphabet = b"abcdefghijklmnopqrstuvwxyz0123456789";
    let mut body = Vec::with_capacity(size);
    for i in 0..size {
        body.push(alphabet[i % alphabet.len()]);
    }
    body
}

fn serialize_body_base64<S>(body: &[u8], serializer: S) -> Result<S::Ok, S::Error>
where
    S: Serializer,
{
    serializer.serialize_str(&STANDARD.encode(body))
}

fn deserialize_body_base64<'de, D>(deserializer: D) -> Result<Vec<u8>, D::Error>
where
    D: Deserializer<'de>,
{
    let encoded = String::deserialize(deserializer)?;
    STANDARD.decode(encoded).map_err(D::Error::custom)
}
