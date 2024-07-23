[![Go](https://img.shields.io/badge/go-1.22.4-blue.svg)](https://golang.org/)
[![Go](https://img.shields.io/badge/version-1.0.0-purple.svg)](https://golang.org/)

## Redis client written in go

### Modes of usage
1. Use as cli tool
2. Use it with your golang project



You can use this on cli but that has limitations. For example, there is not really a way to tell if the value is string, number or array from cli. So, everything is stored as string. We could implement something intelligent but this is unnecessarily complex and unusable and not even official cli tool does it afaik.

## How to use secure tcp
Note: This will only work with SAN certificate and not CN ones as official `redis-server` requires. So make sure to have a configuration file for openssl as per your organization.
For example, you can use the following content as SAN configuration file:
```bash
    [ req ]
    default_bits       = 2048
    distinguished_name = req_distinguished_name
    req_extensions     = req_ext
    x509_extensions    = v3_ca
    prompt             = no

    [ req_distinguished_name ]
    C  = US
    ST = California
    L  = San Francisco
    O  = MyCompany
    OU = MyDivision
    CN = localhost

    [ req_ext ]
    subjectAltName = @alt_names

    [ v3_ca ]
    subjectAltName = @alt_names
    keyUsage = keyCertSign, cRLSign, digitalSignature, keyEncipherment
    extendedKeyUsage = serverAuth, clientAuth

    [ alt_names ]
    DNS.1 = localhost
```

1. Generate a private key
```bash
openssl genpkey -algorithm RSA -out redis.key
```
2. Generate a self signed certificate with SAN using config file 
```bash
openssl req -x509 -nodes -new -key redis.key -sha256 -days 3650 -out redis.crt -config openssl.cnf
```

3. Update redis-server config file with below content to use certificates with SAN then restart redis
```bash
    tls-port 6379
    tls-ca-cert-file /etc/pki/tls/certs/redis.crt
    tls-cert-file /etc/pki/tls/certs/redis.crt
    tls-key-file /etc/pki/tls/certs/redis.key
    tls-auth-clients no
```

4. Use certificate path while connecting to tls configured server. Check any test file e.g pkg/api/del_test.go
```bash
	tlsConfig, err := utils.PrepareTLSConfig("../../redis.crt", "../../redis.key")
```

5. Run tests to check if it works

## TODOS

- [x] Implement pretty printer
- [x] Change string type value to []byte
- [x] TODO: Use \_ instead of ignoring
- [x] Use native integer instead of string
- [x] Expand as golang library
- [x] Support for blocking send due to connection pooling
- [x] Support for secure tcp
