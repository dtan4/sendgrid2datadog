FROM alpine:3.5

COPY bin/sendgrid2datadog /sendgrid2datadog

ENTRYPOINT ["/sendgrid2datadog"]
