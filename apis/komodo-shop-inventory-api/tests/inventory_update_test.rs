// Reserve / release / confirm / restock integration tests.
// Requires spawn_app() to be wired before enabling.

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