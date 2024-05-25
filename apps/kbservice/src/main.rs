pub mod adapters;
pub mod application;
pub mod errors;
pub mod kbs;
pub mod types;

#[tokio::main]
async fn main() {
    println!("⏱️\tStarting kb api application...");
    application::app::run().await;
}
