use axum::{
    http::StatusCode,
    response::{IntoResponse, Response},
    Json,
};
use serde_json::json;
use thiserror::Error;

/// RFC 7807 Problem+JSON error response.
/// Error codes use the 50xxx range (RangePayment = 50) defined in
/// `komodo-forge-sdk-go/http/errors/ranges.go` — the single source of truth
/// for cross-service error code range assignments.
#[derive(Debug, Error)]
pub enum AppError {
    // 50001 — payment-specific
    #[error("Insufficient funds")]
    InsufficientFunds,
    #[error("Payment declined")]
    Declined,
    #[error("Payment method invalid")]
    MethodInvalid,
    #[error("Transaction failed")]
    TransactionFailed,
    #[error("Refund failed")]
    RefundFailed,
    #[error("Payment provider error: {0}")]
    ProviderError(String),

    // General
    #[error("Validation error: {0}")]
    Validation(String),
    #[error("Unauthorized")]
    Unauthorized,
    #[error("Not found")]
    NotFound,
    #[error("Internal server error")]
    Internal(#[from] anyhow::Error),
}

impl AppError {
    fn status(&self) -> StatusCode {
        match self {
            AppError::InsufficientFunds | AppError::Declined => StatusCode::PAYMENT_REQUIRED,
            AppError::MethodInvalid | AppError::Validation(_) => StatusCode::BAD_REQUEST,
            AppError::Unauthorized => StatusCode::UNAUTHORIZED,
            AppError::NotFound => StatusCode::NOT_FOUND,
            AppError::ProviderError(_) => StatusCode::BAD_GATEWAY,
            AppError::TransactionFailed | AppError::RefundFailed | AppError::Internal(_) => {
                StatusCode::INTERNAL_SERVER_ERROR
            }
        }
    }

    fn error_code(&self) -> u32 {
        // 50xxx range
        match self {
            AppError::InsufficientFunds => 50001,
            AppError::Declined => 50002,
            AppError::MethodInvalid => 50003,
            AppError::TransactionFailed => 50004,
            AppError::RefundFailed => 50005,
            AppError::ProviderError(_) => 50006,
            AppError::Validation(_) => 50400,
            AppError::Unauthorized => 50401,
            AppError::NotFound => 50404,
            AppError::Internal(_) => 50500,
        }
    }
}

impl IntoResponse for AppError {
    fn into_response(self) -> Response {
        let status = self.status();
        let body = json!({
            "type": format!("https://komodo.shop/errors/{}", self.error_code()),
            "title": self.to_string(),
            "status": status.as_u16(),
        });
        (status, Json(body)).into_response()
    }
}
