// GET /stock/:sku and GET /stock?skus=... integration tests.
// Requires spawn_app() to be wired before enabling.

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