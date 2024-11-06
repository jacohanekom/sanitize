FROM alpine:latest

WORKDIR /sanitize
COPY sanitize .
COPY sql_sensitive_list.json .

ENTRYPOINT [ "./sanitize" ]
EXPOSE 8080