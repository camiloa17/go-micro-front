# # base go image
# FROM golang:1.20.3-alpine as builder

# RUN mkdir /app

# COPY . /app
# WORKDIR /app
# RUN echo $(ls)

# RUN CGO_ENABLE=0 go build -o brokerApp ./cmd/api

# RUN chmod +x /app/brokerApp

FROM alpine:latest

RUN mkdir /app


#COPY --from=builder /app/brokerApp /app

COPY brokerApp /app

CMD ["/app/brokerApp"]