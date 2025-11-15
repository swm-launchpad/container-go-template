FROM python:3.12-slim AS builder

WORKDIR /app

RUN pip install --no-cache-dir uv
COPY pyproject.toml uv.lock* ./
RUN uv pip install --system --no-cache-dir . && \
    uv pip install --system --no-cache-dir uvicorn

FROM python:3.12-slim

WORKDIR /app

RUN addgroup --system --gid 1001 app && \
    adduser --system --uid 1001 --gid 1001 app

COPY --from=builder /usr/local/lib /usr/local/lib
COPY --from=builder /usr/local/bin /usr/local/bin

COPY . .

RUN chown -R app:app /app

USER app

EXPOSE 8000

CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000", "--workers", "4"]
