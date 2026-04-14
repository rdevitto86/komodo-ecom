// GET /stock/:sku and GET /stock?skus=... integration tests.
// Requires spawn_app() to be wired before enabling.

#[cfg(test)]
mod read_tests {
    use super::*;

    #[tokio::test]
    #[ignore = "requires spawn_app to be implemented"]
    async fn get_stock_returns_stock_level() {
        // let addr = common::spawn_app().await;
        // seed DynamoDB with a test SKU
        // GET /stock/TEST-SKU → 200 StockLevel
        todo!()
    }

    #[tokio::test]
    #[ignore = "requires spawn_app to be implemented"]
    async fn get_stock_unknown_sku_returns_404() {
        todo!()
    }

    #[tokio::test]
    #[ignore = "requires spawn_app to be implemented"]
    async fn batch_stock_returns_map() {
        todo!()
    }
}

// Reserve / release / confirm / restock integration tests.
// Requires spawn_app() to be wired before enabling.

#[cfg(test)]
mod update_tests {
    use super::*;

    #[tokio::test]
    #[ignore = "requires spawn_app to be implemented"]
    async fn reserve_decrements_available_qty() {
        todo!()
    }

    #[tokio::test]
    #[ignore = "requires spawn_app to be implemented"]
    async fn reserve_insufficient_stock_returns_409() {
        todo!()
    }

    #[tokio::test]
    #[ignore = "requires spawn_app to be implemented"]
    async fn release_hold_restores_available_qty() {
        todo!()
    }

    #[tokio::test]
    #[ignore = "requires spawn_app to be implemented"]
    async fn confirm_decrements_reserved_increments_committed() {
        todo!()
    }

    #[tokio::test]
    #[ignore = "requires spawn_app to be implemented"]
    async fn restock_increments_available_qty() {
        todo!()
    }
}