#!/bin/bash
app=online-checker
image=online-checker:0.0.1
current_path=$('pwd')
docker stop $app
docker rm $app
docker rmi $image

go build
docker build -t $image .
docker run --restart=always --name="$app" -d -p 8082:8082 -v $current_path/public:/online-checker/public $image