# acme-dns-certbot-hook

## Overview [![Go Report Card](https://goreportcard.com/badge/github.com/koesie10/acme-dns-certbot-hook)](https://goreportcard.com/report/github.com/koesie10/acme-dns-certbot-hook)

A [Certbot](https://certbot.eff.org/) hook for [acme-dns](https://github.com/joohoi/acme-dns).

## Install

```bash
go install github.com/koesie10/acme-dns-certbot-hook@latest
```

## Usage

This project is for use with the certbot manual plugin. It needs to be run as the `--manual-auth-hook` in the following
manner:

```bash
certbot certonly --manual --manual-auth-hook '/etc/acme-dns/acme-dns-certbot-hook -config /etc/acme-dns/acme_dns.json' --preferred-challenges dns -d example.org
```

In this command, you will need to change `/etc/acme-dns` to the path where you have placed `acme-dns-certbot-hook`
and your config file. You will also need to change the domain and make sure you have set up the domain using
[acme-dns](https://github.com/joohoi/acme-dns).

## Configuration

A sample configuration file:

```json
{
    "acme_dns_url": "https://auth.acme-dns.io",
    "propagation_duration": "10s",
    "domains": {
        "example.org": {
            "allowfrom": [
                "192.168.100.1/24",
                "1.2.3.4/32",
                "2002:c0a8:2a00::0/40"
            ],
            "fulldomain": "8e5700ea-a4bf-41c7-8a77-e990661dcc6a.auth.acme-dns.io",
            "password": "htB9mR9DYgcu9bX_afHF62erXaH2TS7bg9KW3F7Z",
            "subdomain": "8e5700ea-a4bf-41c7-8a77-e990661dcc6a",
            "username": "c36f50e8-4632-44f0-83fe-e070fef28a10"
        }
    }
}
```

The `acme_dns_url` and `propagation_duration` can be overwritten per domain by specifying them along with the other 
information.

The information in the domain section can be directly copied from the response of the `/register` endpoint
of the acme-dns server. Only the `username` and `password` are strictly required.

## License

MIT.
