use tokio::time::{sleep, Duration};

// Новая формула расчета:
// noise_reduction_db = clamp(15 + 0.09*room_area + 1.0*thickness_cm + material_coeff, 20..70)
// material_coeff зависит от материалов перегородки
pub async fn calculate(room_area: f64, thickness_cm: f64, materials: &[String]) -> f64 {
    sleep(Duration::from_secs(5)).await;
    if !room_area.is_finite() || room_area <= 0.0 {
        return 0.0;
    }
    let thickness = if thickness_cm.is_finite() { thickness_cm } else { 0.0 };
    let base = 15.0 + 0.09 * room_area + 1.0 * thickness + material_bonus(materials);
    let clamped = base.max(20.0).min(70.0);
    (clamped * 10.0).round() / 10.0
}

fn material_bonus(materials: &[String]) -> f64 {
    if materials.is_empty() {
        return 6.0;
    }
    let mut sum = 0.0;
    let mut count = 0.0;
    for m in materials {
        sum += material_score(m);
        count += 1.0;
    }
    (sum / count).min(18.0)
}

fn material_score(material: &str) -> f64 {
    let m = material.to_lowercase();
    if m.contains("кирпич") || m.contains("brick") {
        16.0
    } else if m.contains("сэндвич") || m.contains("sandwich") {
        13.0
    } else if m.contains("минват") || m.contains("минеральн") || m.contains("wool") {
        12.0
    } else if m.contains("газобет") || m.contains("aerated") {
        11.0
    } else if m.contains("бетон") || m.contains("concrete") {
        10.0
    } else if m.contains("гипс") || m.contains("gypsum") {
        9.0
    } else if m.contains("дерев") || m.contains("wood") {
        7.0
    } else {
        8.0
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[tokio::test]
    async fn test_calculate() {
        let res = calculate(10.0, 12.0, &[]).await;
        assert!(res >= 0.0);
    }
}
