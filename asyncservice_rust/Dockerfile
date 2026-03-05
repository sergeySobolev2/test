FROM rust:1.81-slim AS builder
WORKDIR /app
COPY Cargo.toml .
COPY src ./src
RUN cargo build --release

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/target/release/asyncservice_rust /app/asyncservice_rust
COPY .env /app/.env
EXPOSE 8000
CMD ["/app/asyncservice_rust"]
