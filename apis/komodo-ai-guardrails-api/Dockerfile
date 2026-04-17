FROM python:3.12-slim
WORKDIR /app
RUN pip install --no-cache-dir pipenv
COPY Pipfile Pipfile.lock* ./
RUN pipenv install --system --deploy 2>/dev/null || pipenv install --system
COPY . .
ENV PORT=7113
EXPOSE 7113
CMD ["python", "-m", "uvicorn", "main:app", "--host", "0.0.0.0", "--port", "7113"]
