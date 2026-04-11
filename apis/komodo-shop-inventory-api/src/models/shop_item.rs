use serde::{Deserialize, Serialize};

/// Minimal reference struct — full catalog shape lives in komodo-shop-items-api.
/// Used here only when inventory lookups need item metadata (e.g. restock_threshold).
#[derive(Debug, Serialize, Deserialize)]
pub struct ShopItemRef {
    pub sku: String,
    pub name: String,
    pub restock_threshold: Option<i64>,
}