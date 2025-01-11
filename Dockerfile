FROM scratch
COPY wait-for /wait-for
COPY extras/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["/wait-for"]
