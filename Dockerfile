FROM bash:5.2
COPY wait-for /usr/local/bin/wait-for
ENTRYPOINT ["wait-for"]
