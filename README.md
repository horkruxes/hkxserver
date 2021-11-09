# HORKRUXES

A decentralized social network based on a distributed authentication system.

## Installation

Unlike Mastodont, the horkrux installation on a linux server is **DEAD SIMPLE**.

### Basic installation, 30s chrono

```bash
mkdir horkruxes && cd horkruxes
wget https://github.com/horkruxes/hkxserver/releases/latest/download/hkxserver_0.7.0_linux_amd64.tar.gz
tar -xzvf horkruxes_xxx.yyy.zzz_linux_amd64.tar.gz #insert version or use auto-completion
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
