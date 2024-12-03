FROM krogertechnology-docker-remote.jfrog.io/alpine:latest

RUN apk add --update --no-cache bind-tools curl vim traceroute mtr tcptraceroute tcpdump net-tools

WORKDIR /app

COPY out/apex-load-generator /app/apex-load-generator

EXPOSE 8080

CMD ["/app/apex-load-generator"]