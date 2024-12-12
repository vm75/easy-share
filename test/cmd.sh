#!/bin/sh

TEST_DIR=/src/test
CONFIG_DIR=${TEST_DIR}/config
VAR_DIR=${TEST_DIR}/var

setup() {
  mkdir -p ${VAR_DIR}
  ln -s ${TEST_DIR}/cmd.sh /bin/cmd
  apk --no-cache update
  apk --no-cache upgrade
  apk --no-cache --no-progress add samba nfs-utils
  addgroup -S smb
  adduser -S -D -H -h /tmp -s /sbin/nologin -G smb -g 'Samba User' smbuser
  mkdir -p /var/lib/nfs/rpc_pipefs /var/lib/nfs/v4recovery
  echo "rpc_pipefs    /var/lib/nfs/rpc_pipefs rpc_pipefs      defaults        0       0" >> /etc/fstab
  echo "nfsd  /proc/fs/nfsd   nfsd    defaults        0       0" >> /etc/fstab
  cp -a /etc/resolv.conf /etc/resolv.conf.bak
}

main() {
  case "$1" in
    run|start)
      if $(command -v docker &> /dev/null) ; then
        sudo docker compose up -d
        sudo docker exec test /src/test/cmd.sh setup
        sudo docker exec -ti test sh
      fi
      ;;
    stop)
      if $(command -v docker &> /dev/null); then
        sudo docker compose down
      fi
      ;;
    setup)
      setup
      ;;
    build)
      cd /src/server
      CGO_ENABLED=1 go build -o ${VAR_DIR}/easy-share .
      ;;
    ports)
      netstat -tulpn | grep LISTEN
      ;;
    *)
      echo "Unknown command: $1"
      exit 1
      ;;
  esac
}

script=$(basename "$0")
cd $(dirname $(readlink -f "$0"))
case "$script" in
  cmd|cmd.sh)
    main "$@"
    ;;
esac
exit 0