pub mod adapters;
pub mod application;
pub mod errors;
pub mod kbs;
pub mod types;

pub async fn run() {
    println!("⏱️\tStarting kb api application...");
    application::app::run().await;
}