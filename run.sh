#!/bin/bash
ACCOUNT=${ACCOUNT:-}
TOKEN=${TOKEN:-}
/usr/local/bin/slacker-pagerduty -account $ACCOUNT -token $TOKEN
