# About

This package defines the tool to secure updates using [The Update Framework](https://theupdateframework.github.io/)

# How to use
0. Set up.
The Tuf tool requires you to set up
   a. Google Cloud Storage

   b. Key Management Service

1. Generate Secrets
The first step of securing your updates is to generate you public-private key pairs.
The tuf tool supports generating Elliptical Curve DSA. You can generate your secrets as follows.
```
bazel run tuf:tuf -- generate-secret --file /tmp/secret_key.json

```
Please do not upload to github or share raw secret key.json file.
If you wish to use your own secrets, Please create a json file with following fields.
```
{"PrivateKey":"-----BEGIN PRIVATE KEY-----\nMHcCAQEEIH4gyo6eaWqnwO+YsurNXFfe0Rqh5mozLIvI4lXz/YVdoAoGCCqGSM49\nAwEHoUQDQgAExBEiFrsujWB8x++q2VtV25IpAIcp/Vx7r3FuFeUU9C1i+EL8QAWl\nY99L3oaODeTMqcDaf4MMD8Iodr6bMI7Mvw==\n-----END PRIVATE KEY-----\n",
"PublicKey":"-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAExBEiFrsujWB8x++q2VtV25IpAIcp\n/Vx7r3FuFeUU9C1i+EL8QAWlY99L3oaODeTMqcDaf4MMD8Iodr6bMI7Mvw==\n-----END PUBLIC KEY-----\n",
"KeyType":"ECDSA256"}
```

2. Upload Secrets.
Once you have generated a public-private key pair, you can use the upload-secrets
command to encrypt this secret and upload them to a private Google Cloud Storage bucket.
In future, the upload-secret command will also generate the Tuf Metadata files and upload them to
a public Google Cloud Storage Bucket.
```
bazel run tuf:tuf -- upload-secrets --root-file /tmp/secret_key.json
```


