#!/bin/bash

# timestamp function
timestamp() {
    date +"%T"
}

#  tmep file for stderr redirects
tmpfile=$(mktemp)

# go build
build () {
    echo "⏲️         $(timestamp): started build script..."
    echo "🏗️         $(timestamp): building..."
    go build 2>tmpfile
    if [ -s tmpfile ]; then
        cat tmpfile
        echo "❌        $(timestamp): compliation error, exiting"
        exit 1
    fi
}

# build docker image
buildDocker() {
    echo "🐋        $(timestamp): building image exmaple:test"
    docker build -t example:test .
}

#  minikube deploy
deploy() {
    echo "🌧️        $(timestamp): delploying to minikube"
    kubectl delete deployment example
    kubectl delete service exmaple
    kubectl apply -f deploy.yml
}

# orchestrate
echo "🤖    Welcome to The Builder v0.1."
if [[ $1 = "build" ]]; then
    if [[ $2 = "docker" ]]; then
      if [[ $3 = "deploy" ]]; then
        build
        buildDocker
        deploy
      else
        build
        buildDocker
        fi
        echo "✔️        $(timestamp): complete"
        echo "👋        $(timestamp): exiting..."
    elif [[ $2 = "bin" ]]; then
        build
        echo "✔️        $(timestamp): complete"
        echo "👋        $(timestamp): exiting..."
    else
        echo "🤔        $(timestamp): missing build argument"
    fi
else
    if [[ $1 == "--help" ]]; then
        echo "build - start a build to produce artifacts"
        echo "docker - produces docker images"
        echo "bin - produces executable binaries"
    else
        echo "🤔        $(timestamp): no arguments passed, type --help for a list of arguments"
    fi
fi
rm -f tmpfile
