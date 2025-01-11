# How to update the `ca-certificates.crt`

Run the following command at the project's root:

```bash
docker run -v "$(pwd)/extras:/extras" -it ubuntu bash -c \
    "apt update && apt install ca-certificates -y && cp /etc/ssl/certs/ca-certificates.crt /extras/."
```

This will copy the `ca-certificates.crt` file to the `extras` directory.
