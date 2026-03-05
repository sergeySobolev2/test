use actix_web::{App, HttpServer};
use dotenv::dotenv;
use serde::Deserialize;
use std::env;

#[derive(Deserialize, Clone)]
pub struct Settings {
    pub host: String,
    pub port: u16,
}

impl Settings {
    pub fn new() -> Self {
        dotenv().ok();
        Self {
            host: env::var("HOST").unwrap_or_else(|_| "0.0.0.0".to_string()),
            port: env::var("LISTEN_PORT").unwrap_or_else(|_| "8000".to_string()).parse().unwrap(),
        }
    }
}

pub async fn start_server(settings: Settings) -> std::io::Result<()> {
    HttpServer::new(move || App::new().configure(crate::handlers::init))
        .bind((settings.host.clone(), settings.port))?
        .run()
        .await
}
