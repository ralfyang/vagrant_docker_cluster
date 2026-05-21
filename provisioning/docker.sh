#!/bin/bash
set -euo pipefail

TARGET_USER="${1:-vagrant}"

export LANG=en_US.UTF-8
export DEBIAN_FRONTEND=noninteractive
export NEEDRESTART_MODE=a

echo ">> Configure needrestart"
mkdir -p /etc/needrestart/conf.d
cat > /etc/needrestart/conf.d/99-vagrant.conf <<'EOF'
$nrconf{restart} = 'a';
$nrconf{kernelhints} = 0;
EOF

echo ">> Install base packages"
apt-get update -y
apt-get install -y \
  ca-certificates \
  curl \
  wget \
  net-tools \
  gnupg

echo ">> Copy SSH files"
mkdir -p /root/.ssh
chmod 700 /root/.ssh

if [ -d /tmp/conf.d ]; then
  find /tmp/conf.d -maxdepth 1 -type f \( \
    -name "id_*" -o \
    -name "*.pub" -o \
    -name "config" -o \
    -name "known_hosts" \
  \) -exec cp {} /root/.ssh/ \; 2>/dev/null || true

  chmod 600 /root/.ssh/* 2>/dev/null || true
  chmod 644 /root/.ssh/*.pub 2>/dev/null || true
fi

echo ">> Configure shell profile"
if id "${TARGET_USER}" >/dev/null 2>&1; then
  USER_HOME="$(getent passwd "${TARGET_USER}" | cut -d: -f6)"

  if [ -n "${USER_HOME}" ] && [ -d "${USER_HOME}" ]; then
    touch "${USER_HOME}/.bash_profile"

    if ! grep -q "alias ll='ls -lia'" "${USER_HOME}/.bash_profile"; then
      echo "alias ll='ls -lia'" >> "${USER_HOME}/.bash_profile"
    fi

    chown "${TARGET_USER}:${TARGET_USER}" "${USER_HOME}/.bash_profile"
  fi
else
  echo "Target user does not exist: ${TARGET_USER}"
  exit 1
fi

echo ">> Configure Docker apt repository"
install -m 0755 -d /etc/apt/keyrings

if [ ! -f /etc/apt/keyrings/docker.asc ]; then
  curl -fsSL https://download.docker.com/linux/ubuntu/gpg \
    -o /etc/apt/keyrings/docker.asc
  chmod a+r /etc/apt/keyrings/docker.asc
fi

. /etc/os-release

echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu ${VERSION_CODENAME} stable" \
  > /etc/apt/sources.list.d/docker.list

echo ">> Install Docker"
apt-get update -y
apt-get install -y \
  docker-ce \
  docker-ce-cli \
  containerd.io \
  docker-buildx-plugin \
  docker-compose-plugin \
  docker-ce-rootless-extras

echo ">> Enable Docker"
systemctl enable docker
systemctl restart docker

echo ">> Add ${TARGET_USER} to docker group"
usermod -aG docker "${TARGET_USER}"

echo ">> Configure passwordless sudo"
echo "${TARGET_USER} ALL=(ALL) NOPASSWD:ALL" > "/etc/sudoers.d/${TARGET_USER}"
chmod 440 "/etc/sudoers.d/${TARGET_USER}"

echo ">> Docker check"
docker --version
docker compose version
systemctl is-active docker
