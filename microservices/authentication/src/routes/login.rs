use axum::{Json, extract::State};
use sqlx::PgPool;

use crate::{
    models::{LoginRequest, User},
    security::{jwt::generate_jwt, password::verify_password},
};

pub async fn login(
    State(pool): State<PgPool>,
    Json(payload): Json<LoginRequest>,
) -> Result<Json<String>, String> {
    let user: User = sqlx::query_as("SELECT * FROM users WHERE email = $1 AND active = true")
        .bind(&payload.email)
        .fetch_one(&pool)
        .await
        .map_err(|_| "Invalid credentials")?;

    if !verify_password(&user.password_hash, &payload.password) {
        return Err("Invalid credentials".into());
    }

    let token = generate_jwt(&user.id.to_string(), "secret");
    Ok(Json(token))
}
