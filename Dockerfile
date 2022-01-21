FROM scratch
COPY wait-for /wait-for
ENTRYPOINT ["/wait-for"]
