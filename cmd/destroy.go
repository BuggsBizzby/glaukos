/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
    "os"
    "log"
	"fmt"
    "os/exec"
    "io/ioutil"
    "strings"
	"github.com/spf13/cobra"
)

var listEnvironments bool

// destroyCmd represents the destroy command
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("destroy called")

        // If list flag is provided, list the environments
        if listEnvironments {
            listAvailableEnvironments()
            return
        }

        // Check for name of environment to destroy or 'all'
        if len(args) == 0 {
            fmt.Println("Please specify an environment name or use 'all' to remove all environments")
            return
        }

        envName := args[0]

        // Remove all environments
        if envName == "all" {
            destroyAllEnvironments()
        } else {
            // Remove specific environment
            destroySpecificEnvironment(envName)
        }

        // Update Caddyfile
        removeCaddyfileRoute(envName)
        
        fmt.Println("Destruction complete.")
	},
}

func init() {
    rootCmd.AddCommand(destroyCmd)
    destroyCmd.Flags().BoolVarP(&listEnvironments, "list", "l", false, "List available environments")

}

func listAvailableEnvironments() {
    // Defining path to environment directories
    envPath := "./docker/configs"

    // Read directories
    dirs, err := ioutil.ReadDir(envPath)
    if err != nil {
        log.Fatalf("Failed to read environments: %s", err)
    }

    // List environments
    for _, d := range dirs {
        if d.IsDir() {
            fmt.Println(d.Name())
        }
    }
}

func destroyAllEnvironments() {
    // Defining path to environment directories
    envDirPath := "./docker/configs"
    envDirs, err := os.ReadDir(envDirPath)
    if err != nil {
        log.Printf("Error reading environments directory: %s\n", err)
        return
    }
    // Loop through directories and tear down environments
    for _, envDir := range envDirs {
        if envDir.IsDir() {
            envName := envDir.Name()
            fmt.Printf("Terminating services for environment: %s\n", envName)
            stopCmd := exec.Command("docker-compose", "-f", fmt.Sprintf("%s/%s/docker-compose.yml", envDirPath, envName), "down")
            if err := stopCmd.Run(); err != nil {
                log.Printf("Error stopping containers for environment: %s: %s\n", envName, err)
            }
        }
    }

    // Looping through and removing all directories files for each environment
    for _, envDir := range envDirs {
        if envDir.IsDir() {
            envPath := fmt.Sprintf("%s/%s", envDirPath, envDir.Name())
            fmt.Printf("Removing directories and files for environment: %s\n", envDir.Name())
            if err := os.RemoveAll(envPath); err != nil {
                log.Printf("Error removing directory for environment: %s: %s\n", envDir.Name(), err)
            }
        }
    }

}

func destroySpecificEnvironment(envName string) {
    // Stop specified environment container
    stopCmd := exec.Command("docker-compose", "-f", fmt.Sprintf("./docker/configs/%s/docker-compose.yml", envName), "down")
    if err := stopCmd.Run(); err != nil {
        log.Printf("Error stopping container for environment: %s: %s\n", envName, err)
    }

    // Remove environment directory
    // Defining environment directory
    envPath := fmt.Sprintf("./docker/configs/%s", envName)
    // Removing environment directory and all files
    if err := os.RemoveAll(envPath); err != nil {
        log.Printf("Error removing directory for environment: %s: %s\n", envName, err)
    }
}

func removeCaddyfileRoute(envName string) {
    // Read Caddyfile
    caddyFileContent, err := ioutil.ReadFile("./docker/configs/Caddyfile")
    if err != nil {
        log.Fatalf("Failed to read Caddyfile: %s", err)
    }

    // If all environments are to be destroyed, reset Caddyfile
    if envName == "all" {
        // Define basic Caddyfile content
        originalCaddyfileContent := `{$SITE_ADDRESS} {

log {
    output stdout
    format console
    level DEBUG
}
}`
        // Writing original content to Caddyfile
        err = ioutil.WriteFile("./docker/configs/Caddyfile", []byte(originalCaddyfileContent), 0777)
        if err != nil {
            log.Fatalf("Failed to reset Caddyfile: %s", err)
        }
        return
    }

    // Removing route for specific environment
    // Creating variable 'routeIdentifier' to look for the line where specific route is located
    routeIdentifier := fmt.Sprintf("handle /%s/*", envName)
    startIndex := strings.Index(string(caddyFileContent), routeIdentifier)

    if startIndex != -1 {
        endIndex := strings.Index(string(caddyFileContent[startIndex:]), "\n\n") + startIndex
        if endIndex == -1 {
            endIndex = len(caddyFileContent)
        }

        updatedContent := string(caddyFileContent[:startIndex]) + string(caddyFileContent[endIndex:])
        err = ioutil.WriteFile("./docker/configs/Caddyfile", []byte(updatedContent), 0777)
        if err != nil {
            log.Fatalf("Failed to remove route from Caddyfile: %s", err)
        }
    }
}




