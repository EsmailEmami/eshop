#! /bin/bash
mkdir -p ./oauth-keys

rm -f ./oauth-keys/oauth-private.key
rm -f ./oauth-keys/oauth-public.key

openssl genrsa -out ./oauth-keys/oauth-private.key 4096
openssl rsa -in ./oauth-keys/oauth-private.key -pubout -out ./oauth-keys/oauth-public.key