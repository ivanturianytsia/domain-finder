./build.sh alpine && ./build.sh image && ./build.sh push

APPNAME="domainfinder"
REGISTRY="207.244.96.79:5000"
MACHINE_NAME="nomi-swarm-3m"

HASH=$(git rev-parse HEAD)

gcloud compute ssh $MACHINE_NAME --command="docker service update --image $REGISTRY/$APPNAME:$HASH $APPNAME"
