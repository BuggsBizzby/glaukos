/*
Copyright © 2023 RIPSFORFUN

*/
package cmd

import (
    "glaukos/embed"
	"fmt"
    "os"
    "os/exec"
    "log"
    "strings"
	"github.com/spf13/cobra"
    "io/ioutil"
    "path/filepath"
)

// Constants for Docker Images

const mitmproxyImage = "glaukos_mitmproxy:latest"
const chromiumImage = "glaukos_chromium:latest"

const mitmproxyRepo = "python"
const mitmproxyTag = "3"
const chromiumRepo = "kasmweb/chromium"
const chromiumTag = "1.14.0-rolling"
const caddyCompose = "./Glaukos/docker-compose-caddy.yml"



// summonCmd represents the summon command
var summonCmd = &cobra.Command{
	Use:   "summon",
	Short: "Build the required docker images",
	Long: `Build the necessary docker images for the chromium service, mitmproxy service, and Caddy service.`,

	Run: func(cmd *cobra.Command, args []string) {

        // Create Glaukos directory
        if err := os.MkdirAll("./Glaukos", 0777); err != nil {
            log.Printf("Failed to create Glaukos directory. Error: %s\n", err)
            return
        }

        // Write docker-compose-caddy.yml file to Glaukos directory
        caddyComposePath := "./Glaukos/docker-compose-caddy.yml"
        if err := ioutil.WriteFile(caddyComposePath, []byte(embed.DockerCaddyCompose), 0777); err != nil {
            log.Printf("Failed to write docker-compose-caddy.yml. Error: %s\n", err)
            return
        }

        // Check if the Docker image already exists
        imageExistsCmd := exec.Command("docker", "image", "ls", "{{.Repository}}::{{.Tag}}")
        output, err := imageExistsCmd.Output()
        if err != nil {
            log.Printf("Failed to fetch Docker images. Error: %s\n", err)
            return
        }
        
        // Build mitmproxy image if it doesn't exist
        if !strings.Contains(string(output), mitmproxyImage) {
            log.Println("Building mitmproxy image..")
            err := buildDockerImage(mitmproxyImage, "mitmproxy")
            if err != nil {
                log.Println("Error building mitmproxy image:", err)
                return
            }
            log.Println("Mitmproxy image built successfully!")
        }
    
        // Build the chromium image if it doesn't exist
        if !strings.Contains(string(output), chromiumImage) {
            log.Println("Building chromium image..")
            err := buildDockerImage(chromiumImage, "chromium")
            if err != nil {
                log.Println("Error building chromium image:", err)
                return
            }
            log.Println("Chromium image built successfully!")
        }

        // Update and generate the Caddyfile
        log.Println("Generating content for Caddyfile..")
        if err := generateCaddyfile(); err != nil {
            log.Println("Error creating Caddyfile:", err)
            return  
        }

        // Check if Caddy container is running
        containerRunningCmd := exec.Command("docker-compose", "-f", caddyCompose, "ps", "-q", "caddy")
        output, err = containerRunningCmd.Output()
        if err != nil {
            log.Printf("Failed to fetch running containers. Error %s\n", err)
            return
        }

        // Start Caddy instance if it's not running
        if strings.TrimSpace(string(output)) == "" {
            log.Println("Starting Caddy instance..")
            err := startCaddyInstance()
            if err != nil {
                log.Println("Error starting Caddy instance:", err)
                return
            }
            log.Println("Caddy instance started successfully!")
        }
	},
}

func init() {
	rootCmd.AddCommand(summonCmd)
}

// Used in Run function to build the desired docker images
func buildDockerImage(imageName, targetName string) error {
    // Prepare dockerBuild directory
    err := prepareDockerBuildContext()
    if err != nil {
        log.Println("Error preparing dockerBuild directory", err)
        return err
    }

    cmd := exec.Command("docker", "build", "--file", "./Glaukos/dockerBuild/Dockerfile", "--target", targetName, "-t", imageName, "./Glaukos/dockerBuild")
    output, err := cmd.CombinedOutput()
    if err != nil {
        log.Println(string(output))
    }
    return err
}

// Used in Run function to generate Caddyfile content 
// To avoid the LetsEncrypt rate limit use the following as the content for the Caddyfile when running multiple tests
// {
//  acme_ca https://acme-staging-v02.api.letsencrypt.org/directory
// }
func generateCaddyfile() error {
    // Ensure Glaukos directory exists
    if err := os.MkdirAll("./Glaukos", 0777); err != nil {
        return err
    }

    content := fmt.Sprintf(`
    `)

    return ioutil.WriteFile("./Glaukos/Caddyfile", []byte(content), 0777)
}

// Used in Run function to create the custom docker network and start the Caddy instance
func startCaddyInstance() error {
    // Create custom docker network
    createNetworkCmd := exec.Command("docker", "network", "create", "mynet")
    networkOutput, networkErr := createNetworkCmd.CombinedOutput()
    if networkErr != nil {
        log.Println(string(networkOutput))
    }

    // Build Caddy image
    cmd := exec.Command("docker-compose", "-f", caddyCompose, "up", "-d")
    output, err := cmd.CombinedOutput()
    if err != nil {
        log.Println(string(output))
    }
    return err
}

func prepareDockerBuildContext() error {
    // Define directory path
    dockerBuildContextDir := "./Glaukos/dockerBuild/"

    // Create dockerBuild directory
    if err := os.MkdirAll(dockerBuildContextDir, 0777); err != nil {
        return err
    }

    // Write Dockerfile to dockerBuild directory
    if err := ioutil.WriteFile(filepath.Join(dockerBuildContextDir, "Dockerfile"), []byte(embed.DockerfileContent), 0777); err != nil {
        return err
    }

    // Write vnc_visual_fixes.py to dockerBuild directory
    if err := ioutil.WriteFile(filepath.Join(dockerBuildContextDir, "vnc_visual_fixes.py"), []byte(embed.VNCVisualFixes), 0777); err != nil {
        return err
    }

    // Write favicon.png to dockerBuild directory
    if err := ioutil.WriteFile(filepath.Join(dockerBuildContextDir, "favicon.png"), embed.Favicon, 0777); err != nil {
        return err
    }

    return nil
}
