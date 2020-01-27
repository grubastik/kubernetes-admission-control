#! /bin/sh
set -o errexit

export APP="${1:-webhook-append-label}"
export NAMESPACE="${2:-default}"
export CSR_NAME="${APP}.${NAMESPACE}.svc"

echo "... creating ${app}.key"
openssl genrsa -out ${APP}.key 2048

echo "... creating ${app}.csr"
cat >csr.conf<<EOF
[req]
req_extensions = v3_req
distinguished_name = req_distinguished_name
[req_distinguished_name]
[ v3_req ]
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names
[alt_names]
DNS.1 = ${APP}
DNS.2 = ${APP}.${NAMESPACE}
DNS.3 = ${CSR_NAME}
DNS.4 = ${CSR_NAME}.cluster.local
EOF
echo "openssl req -new -key ${APP}.key -subj \"/CN=${CSR_NAME}\" -out ${APP}.csr -config csr.conf"
openssl req -new -key ${APP}.key -subj "/CN=${CSR_NAME}" -out ${APP}.csr -config csr.conf
