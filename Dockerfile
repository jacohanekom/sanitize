FROM alpine:latest

WORKDIR /sanitize
COPY sanitize .

ENTRYPOINT [ "./sanitize" ]
EXPOSE 8080