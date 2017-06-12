[![Build Status](https://travis-ci.org/davidwalter0/go-qr-generator.svg?branch=master)](https://travis-ci.org/davidwalter0/go-qr-generator)

# A QR code generator written in Golang
Starts an HTTP server (listening on port 8080) that generates QR codes. Once installed and running (see below), the service accepts the following two parameters:
* ```data```: (Required) The (URL encoded) string that should be encoded in the QR code
* ```size```: (Optional) The size of the image (default: 250)

E.g. ```http://your-domain.tld:8080/?data=Hello%2C%20world&size=300```

## Installation
Download the source code and install it using the `go install` command.

Create certificates, one option is letsencrypt.org like
`https://certbot.eff.org`

Set the environment variables to initialize the certificate location
whether to use https and, if so, tls/https host (matching certs) and
port number

```
    export APP_HTTPS=true
    export APP_HOST=your-domain.tld ;
    export APP_PORT=8443 ;
    export APP_CERT=/etc/letsencrypt/live/your-domain.tld/cert.pem ;
    export APP_KEY=/etc/letsencrypt/live/your-domain.tld/privkey.pem ; 

```

Alternatively, use Docker to run the service in a container:

```
docker run -d -p 8080:8080 davidwalter0/go-qr-generator
```

## References
* Barcode Library: https://github.com/boombuler/barcode

## Author

* [Sam Wierema](http://wiere.ma)

---

*Changes*
- altered environment configuration to include certs, host & port
- tested with letsencrypt.orgs certificate generator instructions
  - https://certbot.eff.org
- Add environment variables for host / port to listen on, default
  to 127.0.0.1:8080
- Add CORS header to enable cross container multi-process connect
  e.g. to be paired in a kubernetes POD in a container

```
    export APP_HTTPS=true;
    export APP_HOST=your-domain.tld ;
    export APP_PORT=8443 ;
    export APP_CERT=/etc/letsencrypt/live/your-domain.tld/cert.pem ;
    export APP_KEY=/etc/letsencrypt/live/your-domain.tld/privkey.pem ; 
    sudo -E /usr/local/go/bin/go run main.go

```
