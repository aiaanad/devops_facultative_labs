FROM python:3.12-slim

# 1. Оптимизация Python (не писать .pyc, выводить логи сразу)
ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1

WORKDIR /app

# 2. Переменная для пути к сервису
ARG SERVICE_PATH

# 3. Установка зависимостей (используем кэш)
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# 4. Копируем код конкретного сервиса
COPY ${SERVICE_PATH}/main.py .

# 5. Создаем не-root пользователя
RUN useradd -m appuser && chown -R appuser /app
USER appuser

# 6. Порт
EXPOSE 8000

CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]