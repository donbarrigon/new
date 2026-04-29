use forge::Server;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let mut server = Server::new("127.0.0.1", 8080);
    server.listen().await?;
    Ok(())
}
