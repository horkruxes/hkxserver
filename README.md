# HORKRUXES

A decentralized social network based on a distributed authentication system.

## Installation

Unlike Mastodont, the horkrux installation on a linux server is **DEAD SIMPLE**.

### Basic installation, 30s chrono

```bash
mkdir horkruxes && cd horkruxes
wget https://github.com/horkruxes/hkxserver/releases/latest/download/hkxserver_linux_amd64.tar.gz
tar -xzvf horkruxes_linux_amd64.tar.gz #insert version or use auto-completion
./horkruxes version
./horkruxes
```

### HTTPS with nginx installed

Copy this to `/etc/nginx/sites-available/horkruxes`

```nginx
server {
    server_name your.server.name; # Use your domain name
    location / {
        proxy_pass http://localhost:80; # You can change the port in hkxconfig.toml
    }
}
```

Then run

```bash
sudo ln -s /etc/nginx/sites-available/horkruxes /etc/nginx/sites-enabled/horkruxes
sudo systemctl restart nginx
sudo certbot --nginx #For HTTPS
```

And that's all folks

### What is signed ?

A list of bytes generated from strings with different encodings, in this order:

- The message (utf-8)
- The public key (base64)
- The Displayed Name (utf-8)

## Development

I like to use [air](https://github.com/cosmtrek/air) to run my projects.

```
go generate ./... # Generates Tailwind styles
go run . # Or use `air`
```

## Project structure

- `/api`: the API routes and definitions
- `/client`: a Golang client for Horkruxes. This is the base for the web client and the cli, and must not depend on service.Service (stateless)
- `/docs`: Swagger generated documentation
- `/exceptions`: explicit
- `/model`: models, independant from anything
- `/query`: database-related operations given a provided service.Service
- `/service`: contain some server state (db and server config)
- `/static`: static resources for the web client
- `/templates`: templates for the web client
- `/views`: web client (mostly Go SSR)
