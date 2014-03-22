from stackbrew/debian:jessie
maintainer evan hazlett "<ehazlett@arcus.io>"
run apt-get update
run apt-get install -y ca-certificates
add slacker-pagerduty /usr/local/bin/slacker-pagerduty
add run.sh /usr/local/bin/run
expose 8080
cmd ["/usr/local/bin/run"]
