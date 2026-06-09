FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY arc-linux-amd64 /app/arc
COPY config/phrases.yaml /app/config/phrases.yaml
EXPOSE 9099
ENTRYPOINT ["/app/arc"]
CMD ["-port=9099", "-config=/app/config/phrases.yaml", "-vosk-url=ws://vosk:2700"]
