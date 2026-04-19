pub async fn health() -> &'static str { "OK" }

#[cfg(test)]
mod tests {
    use super::*;

    #[tokio::test]
    async fn health_returns_ok() {
        assert_eq!(health().await, "OK");
    }
}