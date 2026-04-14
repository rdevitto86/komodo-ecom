/// Shared test helpers — spawn a real app on a random port for integration tests.
///
/// Usage:
/// ```
/// let addr = spawn_app().await;
/// let client = reqwest::Client::new();
/// let res = client.get(format!("{}/health", addr)).send().await.unwrap();
/// ```
pub async fn spawn_app() -> String {
    // TODO: wire up a test DynamoDB (localstack) before enabling
    todo!("spawn test server on random port with DynamoInventoryRepo pointing at LocalStack")
}
