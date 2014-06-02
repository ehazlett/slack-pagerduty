#!/bin/bash
ACCOUNT=${ACCOUNT:-}
TOKEN=${TOKEN:-}
/usr/local/bin/slack-pagerduty -account $ACCOUNT -token $TOKEN
