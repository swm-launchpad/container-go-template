FROM python:3.12-slim AS builder

WORKDIR /app

RUN pip install --no-cache-dir uv
COPY pyproject.toml uv.lock* ./
RUN uv pip install --system --no-cache-dir . && \
    uv pip install --system --no-cache-dir gunicorn

FROM python:3.12-slim

WORKDIR /app

COPY --from=builder /usr/local/lib /usr/local/lib
COPY --from=builder /usr/local/bin /usr/local/bin

COPY . .

ENV DJANGO_SETTINGS_MODULE=config.settings
RUN python manage.py collectstatic --noinput || echo "Warning: collectstatic failed, continuing..."

RUN addgroup --system --gid 1001 app && \
    adduser --system --uid 1001 --gid 1001 app && \
    chown -R app:app /app

USER app

EXPOSE 8000

CMD ["gunicorn", "--bind", "0.0.0.0:8000", "--workers", "4", "config.wsgi"]
