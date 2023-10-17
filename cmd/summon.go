/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
    "os/exec"
    "log"
    "strings"
	"github.com/spf13/cobra"
    "io/ioutil"
)

// Constants for Docker Images

const mitmproxyImage = "glaukos_mitmproxy:latest"
const chromiumImage = "glaukos_chromium:latest"

const mitmproxyRepo = "python"
const mitmproxyTag = "3"
const chromiumRepo = "kasmweb/chromium"
const chromiumTag = "1.14.0-rolling"
const caddyCompose = "./docker/docker-compose-caddy.yml"

// Used in Run function to build the desired docker images
func buildDockerImage(imageName, targetName string) error {
    cmd := exec.Command("docker", "build", "--file", "Dockerfile", "--target", targetName, "-t", imageName, "./docker")
    output, err := cmd.CombinedOutput()
    if err != nil {
        log.Println(string(output))
    }
    return err
}

// Used in Run function to generate Caddyfile content and supply it with the user-given siteAddress value
func modifyCaddyfile(siteAddress string) error {
    content := fmt.Sprintf(`%s {

log {
    output stdout
    format console
    level DEBUG
}
}`, siteAddress)

    return ioutil.WriteFile("./docker/configs/Caddyfile", []byte(content), 0777)
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


// summonCmd represents the summon command
var summonCmd = &cobra.Command{
	Use:   "summon",
	Short: "Build the required docker images",
	Long: `Summon Glaukos from the watery depths to endow you with the necessary abilities to conquer the sea. Dude, it just adds the mitmproxy image and the kasmweb/chromium image.`,

	Run: func(cmd *cobra.Command, args []string) {
        siteAddress, _ := cmd.Flags().GetSTring("siteAddress")

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
        if err := modifyCaddyfile(siteAddress); err != nil {
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
    summonCmd.Flags().StringP("siteAddress", "a", "", "The domain used for routing purposes in the Caddy config file. Ex: `sharepoint.evilcorp.com`")
    summonCmd.MarkFlagRequired("siteAddress")
}
