# generate CA key and certificate
openssl genrsa -out ca.key 4096 # create the rsa private key
openssl req -new -x509 -key ca.key -sha256 -days 365 -out ca.crt -config openssl.cnf # create the certificate

# generate server key and certificate signing request (CSR)
openssl genrsa -out server.key 4096 # create the rsa private key
openssl req -new -key server.key -out server.csr -config openssl.cnf # create the certificate signing request

# generate server certificate
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 365 -sha256 -extfile openssl.cnf -extensions v3_req # create the certificate

# generate client (rider) key and CSR
openssl genrsa -out rider.key 4096 # create the rsa private key
openssl req -new -key rider.key -out rider.csr -config rider_openssl.cnf # create the certificate signing request

# generate client (rider) certificate
openssl x509 -req -in rider.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out rider.crt -days 365 -sha256 -extfile rider_openssl.cnf -extensions v3_req # create the certificate

# generate client (driver) key and CSR
openssl genrsa -out driver.key 4096 # create the rsa private key
openssl req -new -key driver.key -out driver.csr -config driver_openssl.cnf # create the certificate signing request

# generate client (driver) certificate
openssl x509 -req -in driver.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out driver.crt -days 365 -sha256 -extfile driver_openssl.cnf -extensions v3_req # create the certificate
