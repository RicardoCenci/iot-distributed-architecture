#!/bin/bash
password="$1"
output_file="$2"

if [ -z "$password" ]; then
	read -sp "Enter password: " password
	echo ""
fi

if [ -z "$output_file" ]; then
	read -p "Enter output file: " output_file
	echo ""
fi

salt=$(openssl rand -hex 4)
hash=$(echo -n "$(echo -n $salt | xxd -r -p)${password}" | openssl dgst -sha256 -binary | xxd -p -c 256)
password_hash=$(echo -n "${salt}${hash}" | xxd -r -p | base64)
echo $password_hash > $output_file
