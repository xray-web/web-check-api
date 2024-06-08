# Stage 1: Build the application
# Copy over go.mod, install dependencies, copy source code, build the app
FROM golang:1.22.4-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/app main.go

# Stage 2: Run the application
# Create non-root user, copy bin from build stage, set perms, expose port, run app
FROM alpine:latest
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
COPY --from=builder /app/bin/app /usr/local/bin/app
RUN chmod +x /usr/local/bin/app
USER appuser
EXPOSE 8080
CMD ["app"]

