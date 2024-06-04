#!/bin/bash

EXECUTABLE_PATH="$(pwd)/vvmanager"
PID_FILE="$(pwd)/vvmanager.pid"
LOG_DIR="$(pwd)"
LOG_FILE="$LOG_DIR/history.log"

rotate_logs() {
    local max_files=4
    for ((i=max_files; i>0; i--)); do
        if [ -f "$LOG_FILE.$i" ]; then
            mv "$LOG_FILE.$i" "$LOG_FILE.$((i+1))"
        fi
    done
    if [ -f "$LOG_FILE" ]; then
        mv "$LOG_FILE" "$LOG_FILE.1"
    fi
}

case "$1" in
  start)
    if [ -f "$PID_FILE" ]; then
      echo "vvmanager is already running."
      exit 1
    fi
    rotate_logs
    nohup "$EXECUTABLE_PATH" run > "$LOG_FILE" 2>&1 &
    echo $! > "$PID_FILE"
    echo "vvmanager started."
    ;;
  stop)
    if [ ! -f "$PID_FILE" ]; then
      echo "vvmanager is not running."
      exit 1
    fi
    kill $(cat "$PID_FILE")
    rm "$PID_FILE"
    echo "vvmanager stopped."
    ;;
  restart)
    "$0" stop
    "$0" start
    ;;
  run)
    "$EXECUTABLE_PATH" run
    ;;
  *)
    echo "Usage: vvmanager {start|stop|restart|run}"
    exit 1
    ;;
esac

