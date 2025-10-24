#!/bin/sh

# Read the secret file and export key-value pairs as environment variables
while IFS='=' read -r key value || [ -n "$key" ]; do
  export "$key"="$value"
done < /run/secrets/REDIS_SECRET

while IFS='=' read -r key value || [ -n "$key" ]; do
  export "$key"="$value"
done < /run/secrets/RDS_MYSQL_SECRET

# run the app
/my-app