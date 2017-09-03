APPNAME="domainfinder"
REGISTRY="207.244.96.79:5000"
MACHINE_NAME="nomi-swarm-3m"

function build_hash {
  echo "package main
var Build = \"${HASH}\"" > app/Build.go
}

function start {
  Y='\033[1;33m'
  NC='\033[0m'
  printf "${Y}\n - $@...\n${NC}"
}
function complete {
  G='\033[0;32m'
  NC='\033[0m'
  printf "${G}\n - Completed: $@.\n${NC}"
}
function build_binary_mac {
  STEPNAME="Building macOS binary"
  start $STEPNAME
  go get -v -d -t ./app
  go test --cover -v ./app/...
  go build -v -o bin/${APPNAME}_mac ./app
  complete $STEPNAME
}
function build_binary_alpine {
  STEPNAME="Building Alpine Linux binary"
  start $STEPNAME
  go get -v -d -t ./app
  go test --cover -v ./app/...
  CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -v -o bin/${APPNAME}_alpine ./app
  complete $STEPNAME
}
function build_binaries {
  build_binary_mac
  build_binary_alpine
}

function build_image {
  STEPNAME="Building Docker image"
  start $STEPNAME
  docker build --tag ${APPNAME}:latest --file ./deploy/Dockerfile .
  complete $STEPNAME
}
function push_image {
  STEPNAME="Push Docker image to Registry"
  start $STEPNAME

  docker tag ${APPNAME}:latest $REGISTRY/${APPNAME}:latest
  docker push $REGISTRY/${APPNAME}:latest

  docker tag ${APPNAME}:latest $REGISTRY/${APPNAME}:$HASH
  docker push $REGISTRY/${APPNAME}:$HASH

  complete $STEPNAME
}

function deploy {
  STEPNAME="Deploy Docker service localy"
  start $STEPNAME

  docker service update \
    --image $REGISTRY/${APPNAME}:$HASH \
    ${APPNAME}

  complete $STEPNAME
}

function gcloud {
  STEPNAME="Deploy Docker service to Google Cloud"
  start $STEPNAME

  gcloud compute ssh $MACHINE_NAME --command="docker service update --image $REGISTRY/$APPNAME:$HASH $APPNAME"

  complete $STEPNAME
}

HASH=$(git rev-parse HEAD)

case $1 in
  mac)
    build_binary_mac
    ;;
  alpine)
    build_binary_alpine
    ;;
  both)
    build_binaries
    ;;
  image)
    build_image
    ;;
  push)
    push_image
    ;;
  deploy)
    deploy
    ;;
  gcloud)
    gcloud
    ;;
  all)
    build_hash
    build_binaries
    build_image
    push_image
    ;;
esac
