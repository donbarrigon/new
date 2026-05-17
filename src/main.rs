
#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    ironforge::server_start().await?;
    return Ok(())
}
