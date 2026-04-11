use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use uuid::Uuid;

#[derive(Debug, Serialize, Deserialize)]
#[serde(rename_all = "snake_case")]
pub enum CardBrand {
    Visa,
    Mastercard,
    Amex,
    Discover,
    Other,
}

#[derive(Debug, Serialize)]
pub struct PaymentMethod {
    pub id: String,
    pub user_id: Uuid,
    pub last4: String,
    pub brand: CardBrand,
    pub exp_month: u8,
    pub exp_year: u16,
    pub is_default: bool,
    pub created_at: DateTime<Utc>,
}

#[derive(Debug, Deserialize)]
pub struct AddMethodRequest {
    pub provider_token: String,
    pub set_default: Option<bool>,
}

// Payment plans (installments)

#[derive(Debug, Serialize, Deserialize)]
#[serde(rename_all = "snake_case")]
pub enum PlanStatus {
    Active,
    Completed,
    Canceled,
    Defaulted,
}

#[derive(Debug, Serialize, Deserialize)]
#[serde(rename_all = "snake_case")]
pub enum InstallmentStatus {
    Scheduled,
    Paid,
    Failed,
    Skipped,
}

#[derive(Debug, Serialize)]
pub struct Installment {
    pub index: u8,
    pub amount_cents: u64,
    pub due_at: DateTime<Utc>,
    pub paid_at: Option<DateTime<Utc>>,
    pub status: InstallmentStatus,
}

#[derive(Debug, Serialize)]
pub struct PaymentPlan {
    pub plan_id: Uuid,
    pub order_id: Uuid,
    pub user_id: Uuid,
    pub payment_method_id: String,
    pub total_amount_cents: u64,
    pub installment_count: u8,
    pub status: PlanStatus,
    pub installments: Vec<Installment>,
    pub created_at: DateTime<Utc>,
}

#[derive(Debug, Deserialize)]
pub struct CreatePlanRequest {
    pub order_id: Uuid,
    pub payment_method_id: String,
    pub total_amount_cents: u64,
    pub installment_count: u8,
    pub first_due_at: DateTime<Utc>,
    pub interval_days: u16,
}