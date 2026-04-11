use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use uuid::Uuid;

#[derive(Debug, Serialize, Deserialize)]
#[serde(rename_all = "snake_case")]
pub enum ChargeStatus {
    Pending,
    Succeeded,
    Failed,
    Canceled,
}

#[derive(Debug, Serialize, Deserialize)]
#[serde(rename_all = "snake_case")]
pub enum RefundStatus {
    Pending,
    Succeeded,
    Failed,
}

#[derive(Debug, Serialize, Deserialize)]
#[serde(rename_all = "snake_case")]
pub enum PaymentStatus {
    Pending,
    Authorized,
    Captured,
    Refunded,
    PartiallyRefunded,
    Failed,
    Canceled,
}

#[derive(Debug, Deserialize)]
pub struct ChargeRequest {
    pub order_id: Uuid,
    pub amount_cents: u64,
    pub currency: String,
    pub payment_method_id: String,
    pub idempotency_key: Uuid,
}

#[derive(Debug, Serialize)]
pub struct ChargeResponse {
    pub charge_id: Uuid,
    pub order_id: Uuid,
    pub amount_cents: u64,
    pub currency: String,
    pub status: ChargeStatus,
    pub provider_charge_id: String,
    pub created_at: DateTime<Utc>,
}

#[derive(Debug, Deserialize)]
pub struct RefundRequest {
    pub charge_id: Uuid,
    pub amount_cents: Option<u64>,
    pub reason: Option<String>,
}

#[derive(Debug, Serialize)]
pub struct RefundResponse {
    pub refund_id: Uuid,
    pub charge_id: Uuid,
    pub amount_cents: u64,
    pub status: RefundStatus,
    pub provider_refund_id: String,
    pub created_at: DateTime<Utc>,
}
