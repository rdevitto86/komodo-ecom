use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use uuid::Uuid;

#[derive(Debug, Serialize)]
pub struct StockLevel {
    pub sku: String,
    pub available_qty: i64,
    pub reserved_qty: i64,
    pub committed_qty: i64,
    pub restock_threshold: Option<i64>,
}

#[derive(Debug, Deserialize)]
pub struct ReserveRequest {
    pub cart_id: Uuid,
    pub quantity: u32,
}

#[derive(Debug, Serialize)]
pub struct HoldResponse {
    pub hold_id: Uuid,
    pub sku: String,
    pub quantity: u32,
    pub expires_at: DateTime<Utc>,
}

#[derive(Debug, Deserialize)]
pub struct ConfirmRequest {
    pub hold_id: Uuid,
}

#[derive(Debug, Serialize, Deserialize)]
#[serde(rename_all = "snake_case")]
pub enum RestockReason {
    PurchaseReceived,
    ReturnProcessed,
    Adjustment,
}

#[derive(Debug, Deserialize)]
pub struct RestockRequest {
    pub quantity: u32,
    pub reason: Option<RestockReason>,
}

/// Response for batch `/stock` endpoint — map of SKU → StockLevel.
pub type BatchStockResponse = std::collections::HashMap<String, StockLevel>;