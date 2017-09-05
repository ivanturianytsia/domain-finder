APPNAME="domainfinder"
REGISTRY="207.244.96.79:5000"
MACHINE_NAME="nomi-swarm-3m"
HASH=$(git rev-parse HEAD)

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

case $1 in
  mac)
    STEPNAME="Building macOS binary"
    COMMAND="go get -v -d -t ./app && \
      go test --cover -v ./app/... && \
      go build -v -o bin/${APPNAME}_mac ./app"
    ;;
  alpine)
    STEPNAME="Building Alpine Linux binary"
    COMMAND="go get -v -d -t ./app && \
      go test --cover -v ./app/... && \
      CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -v -o bin/${APPNAME}_alpine ./app"
    ;;
  image)
    STEPNAME="Building Docker image"
    COMMAND="docker build --tag ${APPNAME}:latest --file ./deploy/Dockerfile ."
    ;;
  push)
    STEPNAME="Push Docker image to Registry"
    COMMAND="docker tag ${APPNAME}:latest $REGISTRY/${APPNAME}:latest && \
      docker push $REGISTRY/${APPNAME}:latest && \
      docker tag ${APPNAME}:latest $REGISTRY/${APPNAME}:$HASH && \
      docker push $REGISTRY/${APPNAME}:$HASH"
    ;;
  update_service)
    STEPNAME="Update Docker service image"
    COMMAND="docker service update \
      --image $REGISTRY/${APPNAME}:$HASH \
      ${APPNAME}"
    ;;
  docker_service)
    STEPNAME="Create Docker service"
    COMMAND="docker service create \
      -p 8000:8000 \
      --name ${APPNAME} \
      --env BHT_APIKEY=3e942cd3d5ffa9400eb5a25d13552c5a \
      --env PUSHER_ID=392830 \
      --env PUSHER_KEY=80e0323979d179949997 \
      --env PUSHER_SECRET=91a77850826ea6d50f52 \
      --env PORT=8000 \
      $REGISTRY/${APPNAME}:$HASH"
    ;;
esac

if [ "$2" = "print" ]
then
  echo $COMMAND
else
  eval $COMMAND
  complete $STEPNAME
fi
