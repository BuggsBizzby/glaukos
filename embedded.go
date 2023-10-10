package main

import _ "embed"

// Not ready yet // go:embed docker/Caddyfile
var CaddyfileTemplate string

//go:embed docker/docker-compose-template.yml
var DockerComposeTemplate string


