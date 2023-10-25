/*
Copyright © 2023 RIPSFORFUN

*/
package cmd

import (
    "os"
    "log"
	"fmt"
    "os/exec"
    "io/ioutil"
    "regexp"
	"github.com/spf13/cobra"
)

var listEnvironments bool
var burnItDownFlag bool

// destroyCmd represents the destroy command
var destroyCmd = &cobra.Command{
	Use:   "destroy [ENVIRONMENT_NAME | all]",
	Short: "A brief description of your command",
	Long: `Destroy will tear down specified environments. 
Use the '--burn-it-down' flag to remove environment directories and associated files.`,
    Example: `
  # Destroy a specific environment:
  destroy envName
  
  # Destroy all environments:
  destroy all
  
  # Remove environment directories and associated files:
  destroy envName --burn-it-down
  
  # Remove directories and files for all environments:
  destroy all --burn-it-down
`,

	Run: func(cmd *cobra.Command, args []string) {

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

        // Check for burn-it-down flag
        burnItDown := burnItDownFlag

        // Remove all environments
        if envName == "all" {
            destroyAllEnvironments(burnItDown)
        } else {
            // Remove specific environment
            destroySpecificEnvironment(envName, burnItDown)
        }

        // Update Caddyfile
        removeCaddyfileRoute(envName)
        
        fmt.Println("Destruction complete.")
	},
}

func init() {
    rootCmd.AddCommand(destroyCmd)
    destroyCmd.Flags().BoolVarP(&listEnvironments, "list", "l", false, "List available environments")
    destroyCmd.Flags().BoolVarP(&burnItDownFlag, "burn-it-down", "b", false, "Remove environment directories and associated files")

}

func listAvailableEnvironments() {
    // Defining path to environment directories
    envPath := "./Glaukos/environments"

    // Read directories
    dirs, err := ioutil.ReadDir(envPath)
    if err != nil {
        log.Fatalf("Failed to read environments: %s", err)
    }

    // List environments
    fmt.Println("Environments:")
    for _, d := range dirs {
        if d.IsDir() {
            fmt.Println(d.Name())
        }
    }
}

func destroyAllEnvironments(burnItDown bool) {
    // Defining path to environment directories
    envDirPath := "./Glaukos/environments"
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
    if burnItDown {
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
}

func destroySpecificEnvironment(envName string, burnItDown bool) {
    // Stop specified environment container
    stopCmd := exec.Command("docker-compose", "-f", fmt.Sprintf("./Glaukos/environments/%s/docker-compose.yml", envName), "down")
    if err := stopCmd.Run(); err != nil {
        log.Printf("Error stopping container for environment: %s: %s\n", envName, err)
    }

    if burnItDown {
        // Remove environment directory
        // Defining environment directory
        envPath := fmt.Sprintf("./Glaukos/environments/%s", envName)
        // Removing environment directory and all files
        if err := os.RemoveAll(envPath); err != nil {
            log.Printf("Error removing directory for environment: %s: %s\n", envName, err)
        }
    }
}

func removeCaddyfileRoute(envName string) {
    // Read Caddyfile
    caddyFileContent, err := ioutil.ReadFile("./Glaukos/Caddyfile")
    if err != nil {
        log.Fatalf("Failed to read Caddyfile: %s", err)
    }

    // If all environments are to be destroyed, reset Caddyfile
    if envName == "all" {
        // Wipe Caddyfile clean
        err = ioutil.WriteFile("./Glaukos/Caddyfile", []byte(""), 0777)
        if err != nil {
            log.Fatalf("Failed to reset Caddyfile: %s", err)
        }
        // Reload Caddy Service
        updateCaddyService()
        return
    }

    // Removing route for specific environment
    routeRegexPattern := fmt.Sprintf(`(?ms)^%s\..*?{.*?log {\s*output stdout\s*format console\s*level DEBUG\s*}\s*}\s*\n`, regexp.QuoteMeta(envName))
    re := regexp.MustCompile(routeRegexPattern)

        updatedContent := re.ReplaceAllString(string(caddyFileContent), "")
        
        err = ioutil.WriteFile("./Glaukos/Caddyfile", []byte(updatedContent), 0777)
        if err != nil {
            log.Fatalf("Failed to remove route from Caddyfile: %s", err)
        }
}

func updateCaddyService() error {
    // Command to take the Caddy service down
    downCmd := exec.Command("docker-compose", "-f", "./Glaukos/docker-compose-caddy.yml", "down")
    downOutput, downErr := downCmd.CombinedOutput()
    if downErr != nil {
        log.Println("Failed to take Caddy service down:", string(downOutput))
        return downErr
    }

    // Command to bring Caddy service back up
    upCmd := exec.Command("docker-compose", "-f", "./Glaukos/docker-compose-caddy.yml", "up", "-d")
    upOutput, upErr := upCmd.CombinedOutput()
    if upErr != nil {
        log.Println("Failed to bring Caddy service up:", string(upOutput))
        return upErr
    }

    log.Println("Caddy reloaded successfully!")
    return nil
}





