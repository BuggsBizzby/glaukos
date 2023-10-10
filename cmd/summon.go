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
)

// Constants for Docker Images

const mitmproxyImage = "python:3"
const chromiumImage = "kasmweb/chromium:1.14.0-rolling"

const mitmproxyRepo = "python"
const mitmproxyTag = "3"
const chromiumRepo = "kasmweb/chromium"
const chromiumTag = "1.14.0-rolling"


// Used in Run function to build the desired docker images
func buildDockerImage(imageName, targetName string) error {
    cmd := exec.Command("docker", "build", "--file", "Dockerfile", "--target", targetName, "-t", imageName, "./docker")
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
		fmt.Println("summon called")

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

	},
}


func init() {
	rootCmd.AddCommand(summonCmd)
}
