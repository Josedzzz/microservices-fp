use axum::{Router, routing::post};
use dotenvy::dotenv;
use tokio::net::TcpListener;

mod config;
mod db;
mod models;
mod routes;
mod security;

use routes::login::login;

#[tokio::main]
async fn main() {
    dotenv().ok();

    let config = config::Config::from_env();
    let pool = db::connect(&config.database_url).await;

    let app = Router::new()
        .route("/auth/login", post(login))
        .with_state(pool);

    let listener = TcpListener::bind("0.0.0.0:8081").await.unwrap();

    axum::serve(listener, app).await.unwrap();
}
