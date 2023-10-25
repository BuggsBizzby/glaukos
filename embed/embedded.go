package embed

import _ "embed"

//go:embed docker-compose-template.yml
var DockerComposeTemplate string

//go:embed docker-compose-caddy.yml
var DockerCaddyCompose string

//go:embed Dockerfile
var DockerfileContent string

//go:embed vnc_visual_fixes.py
var VNCVisualFixes string

//go:embed favicon.png
var Favicon []byte
