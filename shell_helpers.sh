#!/usr/bin/env bash

COLOR_BLACK="\033[0;30m"
COLOR_RED="\033[0;31m"
COLOR_GREEN="\033[0;32m"
COLOR_YELLOW="\033[0;33m"
COLOR_BLUE="\033[0;34m"
COLOR_MAGENTA="\033[0;35m"
COLOR_CYAN="\033[0;36m"
COLOR_WHITE="\033[0;37m"
COLOR_BRIGHT_BLACK="\033[1;30m"
COLOR_BRIGHT_RED="\033[1;31m"
COLOR_BRIGHT_GREEN="\033[1;32m"
COLOR_BRIGHT_YELLOW="\033[1;33m"
COLOR_BRIGHT_BLUE="\033[1;34m"
COLOR_BRIGHT_MAGENTA="\033[1;35m"
COLOR_BRIGHT_CYAN="\033[1;36m"
COLOR_BRIGHT_WHITE="\033[1;37m"
COLOR_RESET="\033[0m"

function _log_date_header() {
    local COLOR1=${1:-${COLOR_BLUE}}
    local COLOR2=${2:-${COLOR_WHITE}}
    echo -en "${COLOR_RESET}${COLOR1}[${COLOR2}$(date)${COLOR1}] ${COLOR_RESET}"
}

function _error() {
    DATE_HEADER=$(_log_date_header $COLOR_BRIGHT_RED $COLOR_BRIGHT_YELLOW)
    echo -e "${DATE_HEADER} ${COLOR_RED}ERROR:${COLOR_RESET} $1" >&2
}

function _log() {
    DATE_HEADER=$(_log_date_header)
    echo -e "${DATE_HEADER} $1${COLOR_RESET}"
}

function _log2() {
    DATE_HEADER=$(_log_date_header)
    echo -e "${DATE_HEADER} $1${COLOR_RESET}" >&2
}

