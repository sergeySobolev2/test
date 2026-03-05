use crate::calculation::calculate;
use crate::models::{AsyncResult, CalcRequest};
use actix_web::{get, post, web, HttpResponse};
use reqwest::Client;
use serde_json::json;
use std::env;
use tokio::spawn;

#[post("/calculate")]
pub async fn handle_request(req: web::Json<CalcRequest>) -> HttpResponse {
    start_calculation(req.into_inner());
    HttpResponse::Accepted().json(json!({"message": "accepted"}))
}

#[get("/")]
pub async fn handle_root() -> HttpResponse {
    HttpResponse::Ok().json(json!({"status": "ok"}))
}

pub fn init(cfg: &mut web::ServiceConfig) {
    cfg.service(handle_request).service(handle_root);
}

fn start_calculation(req: CalcRequest) {
    spawn(async move {
        let thickness = req.thickness_cm.unwrap_or(0.0);
        let materials = req.materials.unwrap_or_default();
        let result = calculate(req.room_area, thickness, &materials).await;
        if let Err(err) = send_result(req.calculation_id, req.room_area, result).await {
            eprintln!("Failed to send result: {err:?}");
        }
    });
}

async fn send_result(calculation_id: u32, room_area: f64, noise_reduction_db: f64) -> Result<(), reqwest::Error> {
    dotenv::dotenv().ok();
    let backend_url = env::var("BACKEND_RESULT_URL")
        .unwrap_or_else(|_| "http://127.0.0.1:8082/api/async/result".to_string());
    let token = env::var("ASYNC_TOKEN")
        .unwrap_or_else(|_| "async-partition-soundproof-key".to_string());

    let client = Client::builder().no_proxy().build()?;
    client
        .post(backend_url)
        .json(&AsyncResult {
            calculation_id,
            noise_reduction_db,
            room_area,
            service_token: token,
        })
        .send()
        .await?
        .error_for_status()?;

    Ok(())
}
