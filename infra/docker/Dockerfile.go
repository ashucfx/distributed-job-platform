FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy all go modules config
COPY go.work ./
COPY apps/ apps/
COPY packages/ packages/
COPY services/ services/

# Build the specified service
ARG SERVICE_PATH
RUN ls -laR .
RUN cd ${SERVICE_PATH} && go build -o /bin/service .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /bin/service /app/service

CMD ["/app/service"]
