wget --no-check-certificate --spider -S $1 2>&1 | awk '/HTTP\// {print $2}'

curl --silent --insecure -i $1 | grep HTTP | awk '{print $2}'

