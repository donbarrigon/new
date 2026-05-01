
#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    forge::server_start().await?;
    return Ok(())
}
