use serde::{Deserialize, Serialize};

#[derive(Deserialize, Clone)]
pub struct CalcRequest {
    pub calculation_id: u32,
    pub room_area: f64,
    pub thickness_cm: Option<f64>,
    pub materials: Option<Vec<String>>,
}

#[derive(Serialize)]
pub struct AsyncResult {
    pub calculation_id: u32,
    pub noise_reduction_db: f64,
    pub room_area: f64,
    pub service_token: String,
}
