#!/bin/bash
app=online-checker
image=online-checker:0.0.1
current_path=$('pwd')
docker stop $app
docker rm $app
docker rmi $image
docker build -t $image .
docker run --name="$app" -d -p 8082:8082 -v $current_path/public:/online-checker/public $image