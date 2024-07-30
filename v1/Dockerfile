FROM golang:1.22.0-alpine3.19
RUN mkdir /app
COPY . /app

WORKDIR /app
RUN go build -o main ./cmd/api 
CMD [ "/app/main" ]