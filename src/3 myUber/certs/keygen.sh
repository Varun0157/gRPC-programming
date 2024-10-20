# Generate CA key and certificate
openssl genrsa -out ca.key 4096
openssl req -new -x509 -key ca.key -sha256 -subj "/C=US/ST=NJ/O=MyUber CA" -days 365 -out ca.crt

# Generate server key and certificate signing request (CSR)
openssl genrsa -out server.key 4096
openssl req -new -key server.key -out server.csr -subj "/C=US/ST=NJ/O=MyUber Server/CN=localhost"

# Generate server certificate
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 365 -sha256

# Generate client (rider) key and CSR
openssl genrsa -out rider.key 4096
openssl req -new -key rider.key -out rider.csr -subj "/C=US/ST=NJ/O=MyUber Rider/CN=localhost"

# Generate client (rider) certificate
openssl x509 -req -in rider.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out rider.crt -days 365 -sha256

# Generate client (driver) key and CSR
openssl genrsa -out driver.key 4096
openssl req -new -key driver.key -out driver.csr -subj "/C=US/ST=NJ/O=MyUber Driver/CN=localhost"

# Generate client (driver) certificate
openssl x509 -req -in driver.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out driver.crt -days 365 -sha256
