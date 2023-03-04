# STEP 1
FROM golang:alpine as builder
RUN apk update && apk add --no-cache git
RUN mkdir /build
ADD . /build
WORKDIR /build
RUN go install 
RUN go build -o cicdexample
# Step 2
FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/ /app/
WORKDIR /app
CMD ["./cicdexample"]
