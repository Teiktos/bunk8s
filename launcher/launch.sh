#!/bin/bash

echo '########### start testing process ###########'

docker run --name bunk8s-launcher  bunk8s-launcher-image:"${CONFIG_FILE_NAME%.*}"
docker logs bunk8s-launcher > ./serverReply.json
docker rm bunk8s-launcher


