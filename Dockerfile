FROM krogertechnology-docker-remote.jfrog.io/alpine:latest

RUN apk add --update --no-cache bind-tools curl vim traceroute mtr tcptraceroute tcpdump net-tools

WORKDIR /app

COPY out/apex-load-generator /app/apex-load-generator
COPY swagger.yaml /app/swagger.yaml

EXPOSE 8080

CMD ["/app/apex-load-generator"]