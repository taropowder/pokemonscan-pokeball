
CONFIG_KEY_DIR="config/key/"

mkdir -p ${CONFIG_KEY_DIR}

# 根证书
openssl genrsa -out ${CONFIG_KEY_DIR}/ca.key 2048
openssl req -new -sha256 -x509 -days 3650 -key ${CONFIG_KEY_DIR}/ca.key -out ${CONFIG_KEY_DIR}/ca.pem -subj "/C=cn/OU=custer/O=custer/CN=pokemon.go"

# 服务端证书
openssl genpkey -algorithm RSA -out ${CONFIG_KEY_DIR}/server.key
openssl req -new -sha256 -nodes -key ${CONFIG_KEY_DIR}/server.key -out ${CONFIG_KEY_DIR}/server.csr -days 3650 -subj "/C=cn/OU=custer/O=custer/CN=pokemon.go" -config ./openssl.cnf -extensions v3_req
openssl x509 -req -sha256 -days 3650 -in ${CONFIG_KEY_DIR}/server.csr -out ${CONFIG_KEY_DIR}/server.pem -CA ${CONFIG_KEY_DIR}/ca.pem -CAkey ${CONFIG_KEY_DIR}/ca.key -CAcreateserial -extfile ./openssl.cnf -extensions v3_req

openssl genpkey -algorithm RSA -out ${CONFIG_KEY_DIR}/client.key
openssl req -new -sha256 -nodes -key ${CONFIG_KEY_DIR}/client.key -out ${CONFIG_KEY_DIR}/client.csr -days 3650 -subj "/C=cn/OU=custer/O=custer/CN=pokemon.go" -config ./openssl.cnf -extensions v3_req
openssl x509 -req -sha256 -days 3650 -in ${CONFIG_KEY_DIR}/client.csr -out ${CONFIG_KEY_DIR}/client.pem -CA ${CONFIG_KEY_DIR}/ca.pem -CAkey ${CONFIG_KEY_DIR}/ca.key -CAcreateserial -extfile ./openssl.cnf -extensions v3_req

docker-compose up -d

docker exec pokemon_daycare /app/daycare --config config/daycare_config.yml init