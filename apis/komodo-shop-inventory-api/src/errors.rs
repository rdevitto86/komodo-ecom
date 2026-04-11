use axum::{
    http::StatusCode,
    response::{IntoResponse, Response},
    Json,
};
use serde_json::json;
use thiserror::Error;

/// RFC 7807 Problem+JSON error response.
/// Error codes mirror the Go pkg/v1/models/errors.go ranges (44xxx).
#[derive(Debug, Error)]
pub enum AppError {
    // 44xxx — inventory-specific
    #[error("Insufficient stock")]
    InsufficientStock,
    #[error("SKU not found")]
    SkuNotFound,
    #[error("Stock hold not found")]
    HoldNotFound,
    #[error("Stock hold expired")]
    HoldExpired,

    // General
    #[error("Validation error: {0}")]
    Validation(String),
    #[error("Unauthorized")]
    Unauthorized,
    #[error("Internal server error")]
    Internal(#[from] anyhow::Error),
}

impl AppError {
    fn status(&self) -> StatusCode {
        match self {
            AppError::InsufficientStock => StatusCode::CONFLICT,
            AppError::SkuNotFound | AppError::HoldNotFound => StatusCode::NOT_FOUND,
            AppError::HoldExpired => StatusCode::GONE,
            AppError::Validation(_) => StatusCode::BAD_REQUEST,
            AppError::Unauthorized => StatusCode::UNAUTHORIZED,
            AppError::Internal(_) => StatusCode::INTERNAL_SERVER_ERROR,
        }
    }

    fn error_code(&self) -> u32 {
        // 44xxx range
        match self {
            AppError::InsufficientStock => 44001,
            AppError::SkuNotFound => 44002,
            AppError::HoldNotFound => 44003,
            AppError::HoldExpired => 44004,
            AppError::Validation(_) => 44400,
            AppError::Unauthorized => 44401,
            AppError::Internal(_) => 44500,
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
