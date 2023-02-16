#!/bin/bash

WORKDIR="/opt/pokeball"
POKEBALL_PATH="/opt/pokeball/bin/pokeball"
POKEBALL_ATTACHMENT_PATH="/opt/pokeball/attachment/pokeball.tar.gz"

cd $WORKDIR

SERVER_ADDRESS=$(cat /opt/pokeball/config/server)
UPDATE_ADDRESS=$(cat /opt/pokeball/config/update_server)

WORKER_CONFIG_ID=""

POKEBALL_IMAGES=(
pokeball_enscan
pokeball_nuclei
pokeball_oneforall
pokeball_rad
pokeball_xray
)


start(){
      echo "start pokeball worker"
      $POKEBALL_PATH run -s "$SERVER_ADDRESS"
}

stop(){
      echo "stop pokeball worker"
}

check_update(){
  echo "check pokeball need update"
  LAST_HASH=$(curl -s "$UPDATE_ADDRESS/pokemon/md5")
  NOW_HASH=$(md5sum $POKEBALL_ATTACHMENT_PATH | awk '{print $1}' )
  echo  "LAST_HASH is $LAST_HASH, NOW_HASH is $NOW_HASH"
    if [ "$LAST_HASH" != "$NOW_HASH" ]; then
      echo "pokeball is updating"
      mv $POKEBALL_ATTACHMENT_PATH "$POKEBALL_ATTACHMENT_PATH.bak"
      curl -o $POKEBALL_ATTACHMENT_PATH "$UPDATE_ADDRESS/pokemon/download"
      tar -zxvf "$WORKDIR/attachment/pokeball.tar.gz" -C $WORKDIR
      systemctl daemon-reload
    else
      echo "pokeball is the last, no need to update"
    fi
}

check_env(){
#  check docker
  if ! docker info > /dev/null ; then
      echo "error for : docker info"
      exit 1
  fi
  echo "check  docker success"

#  check images
  for IMAGE in  "${POKEBALL_IMAGES[@]}"
  do
  echo "check $IMAGE";
    if ! docker images  | grep "$IMAGE" > /dev/null ; then
      echo "error images : no such image '$IMAGE'"
      images_init
    fi
  done
  echo "check  images success"

#  check network
  if ! docker network ls | grep "pokemon_net" > /dev/null ; then
      docker network create --driver=bridge --subnet=192.161.0.0/16 pokemon_net
      echo "create docker network : docker network create --driver=bridge --subnet=192.161.0.0/16 pokemon_net "
  fi
  echo "check docker network success"
}

images_init(){

  for IMAGE in  "${POKEBALL_IMAGES[@]}"
  do
    # pull images
    docker pull "taropowder/$IMAGE"
  done


}


if [ "$1" == "start" ]; then
  check_env
  check_update
  start
fi

if [ "$1" == "stop" ]; then
  stop
fi

if [ "$1" == "check_env" ]; then
  check_env
fi

if [ "$1" == "check_update" ]; then
  check_update
fi

if [ "$1" == "images_init" ]; then
  images_init
fi


echo "success run pokeball worker"
