docker build -t frontend-service .
docker run -d --rm -p 13000:3000  frontend-service