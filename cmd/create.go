/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
    "strings"
	"github.com/spf13/cobra"
    "os"
    "os/exec"
    "log"
    "io/ioutil"
    "path/filepath"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create new environments",
	Long: `Initializes one or more new Docker environments, allowing for rapid setup based on input for number of environments, the domain configuration, and path names associated with each.`,

	Run: func(cmd *cobra.Command, args []string) {
        envCount, _ := cmd.Flags().GetInt("environments")
        prefix, _ := cmd.Flags().GetString("prefix")
        names, _ := cmd.Flags().GetStringSlice("names")
        targetURL, _ := cmd.Flags().GetString("targetURL")
        siteAddress, _ := cmd.Flags().GetString("siteAddress")

        // Check if both the names and prefix arguments are being used, if so, throw an error.
        if len(names) > 0 && len(prefix) > 0 {
            fmt.Println(len(names))
            fmt.Println("Conflicting arguments: You cannot provide both names and a prefix, please choose one.")
            return
        }

        // Check for the presence of a value for names, and if it matches the value provided for the environment count
        if len(names) > 0 {
            if len(names) != envCount {
                fmt.Printf("Mismatch: You provided %d names but requested %d environments.\n", len(names), envCount)
                return
            }
            // If count matches up execute the createEnvironment() with name values 
            for _, name := range names {
                createEnvironment(name, targetURL, siteAddress)
            }
        // Check if prefix is being used and has a value    
        } else if len(prefix) > 0 {
            for i := 1; i <= envCount; i+=1 {
                // Append numbers to the end of the prefix (Ex. --prefix test = test1, test2, etc.)
                envName := fmt.Sprintf("%s%d", prefix, i)
                // Execute createEnvironment() with prefix values 
                createEnvironment(envName, targetURL, siteAddress)
            }
        // Check if no path naming scheme has been provided, if so throw an error.
        } else {
            fmt.Println("Path naming scheme not defined: You must provide either prefix or names.")
            return
        }
	},
}

func init() {
    rootCmd.AddCommand(createCmd)

    // Defining flags
    createCmd.Flags().IntP("environments", "e", 1, "Number of environments to create")
    createCmd.MarkFlagRequired("environments")
    createCmd.Flags().StringP("prefix", "p", "", "Prefix for naming environments")
    createCmd.Flags().StringSliceP("names", "n", []string{}, "Comma-separated list of environment names")
    createCmd.Flags().StringP("targetURL", "u", "", "URL of the target website to be displayed to the victim - Ex. 'https://login.microsoftonline.com'")
    createCmd.MarkFlagRequired("targetURL")
    createCmd.Flags().StringP("siteAddress", "a", "", "The domain, used for routing purposes in the Caddy config file. Ex. 'sharepoint.evilcorp.com'")
    createCmd.MarkFlagRequired("siteAddress")
}

func createEnvironment(name, targetURL, siteAddress string) {
    fmt.Println("Creating environment:", name)

    // Generating docker-compose.yml, from template file
    dockerDir := generateCompose(name)    
    if dockerDir == "" {
        return
    }

    // Generating .env files for each environment
    err := generateEnvFile(name, targetURL, siteAddress)
    if err != nil {
        log.Printf("Failed to create .env file. Error: %s\n", err)
        return
    }

    // Update caddy file with new environment
    errCaddy := updateCaddy(name)
    if errCaddy != nil {
        log.Printf("Failed to update Caddyfile for environment %s. Error: %s\n", errCaddy)
        return
    }

//    // Check if the Docker image already exists
//    imageExistsCmd := exec.Command("docker", "image", "ls", "{{.Repository}}::{{.Tag}}")
//    output, err := imageExistsCmd.Output()
//    if err != nil {
//        log.Printf("Failed to fetch Docker images. Error: %s\n", err)
//        return
//    }
//    
//    // Inform user to run the summon command if images are missing
//    if !(strings.Contains(string(output), mitmproxyRepo) && strings.Contains(string(output), mitmproxyTag)) || !(strings.Contains(string(output), chromiumRepo) && strings.Contains(string(output), chromiumTag)) {
//        log.Println("Required Docker images are missing. Please run the `summon` command before creating environments.") 
//        return
        

    // PREVIOUS CODE
//    if !strings.Contains(string(output), mitmproxyImage) || !strings.Contains(string(output), chromiumImage) {
//        log.Println("Required Docker images are missing. Please run the `summon` command before creating environments.")
//        return
//    } 

    // Run containers using docker-compose, within dockerDir
    composeFilePath := filepath.Join(dockerDir, "docker-compose.yml")
    // DEBUGGING
    fmt.Println(composeFilePath)
    run := exec.Command("docker-compose", "-f", composeFilePath, "up", "-d")
    _, err = run.CombinedOutput()
    if err != nil {
        log.Println("Error starting services:", err)
        return
    }
    log.Printf("Environment %s: Services started successfully!", name)
}

// Generate docker-compose.yml from template file
func generateCompose(envName string) string {
    // Read content from docker-compose template file in ./docker
    composeContent, err := ioutil.ReadFile("./docker/docker-compose-template.yml")
    if err != nil {
        log.Printf("Failed to read docker-compose template. Error %s\n", err)
        return ""
    }
    // Modify docker-compose template with user supplied environment names
    modifiedCompose := strings.ReplaceAll(string(composeContent), "{{env_name}}", envName)
    

    // Create unique directory for new environment's config files
    configDir := fmt.Sprintf("./docker/configs/%s", envName)
    if err := os.MkdirAll(configDir, 0755); err != nil {
        log.Printf("Failed to create config directory. Error: %s\n", err)
        return ""
    }

    // Write modified docker-compose file to environment config directory
    if err := ioutil.WriteFile(filepath.Join(configDir, "docker-compose.yml"), []byte(modifiedCompose), 0777); err != nil {
        log.Printf("Failed to write docker-compose. Error: %s\n", err)
        return ""
    }


    return configDir
}

// Update the Caddyfile
func updateCaddy(envName string) error {
    caddyFilePath := "./docker/configs/Caddyfile"

    // Read contents of Caddyfile template
    caddyContent, err := ioutil.ReadFile(caddyFilePath)
    if err != nil {
        log.Printf("Failed to read Caddyfile. Error: %s\n", err)
        return err
    }

    // Create handle directorive for new environment
    handleDirective := fmt.Sprintf(`
handle /%s/* {
    reverse_proxy * {
        to https://%s_chromium:6901
        header_up Authorization "Basic a2FzbV91c2VyOmFzZGZmZHNh"
        transport http {
            tls
            tls_insecure_skip_verify
        }
    }
}
    `, envName, envName)

    // Insert handle directive above log section
    updatedCaddyContent := strings.Replace(string(caddyContent), "\nlog {", handleDirective+"\n\nlog {", 1)
    
    // Save updated Caddyfile
    err = ioutil.WriteFile(caddyFilePath, []byte(updatedCaddyContent), 0644)
    if err != nil {
        log.Printf("Failed to write updated Caddyfile: %s", err)
        return err
    }

    return nil
}

// Generating .env files for each environment
func generateEnvFile(envName, targetURL, siteAddress string) error {
    // Specifying the contents of the .env file
    content := fmt.Sprintf(
    `TARGET_URL=%s
    SITE_ADDRESS=%s
    LANGUAGE=en
    LANG=en_US.UTF-8
    LC_ALL=en_US.UTF-8`,
    targetURL, siteAddress)
    
    // Assign configDir to the current environment directory
    configDir := fmt.Sprintf("./docker/configs/%s", envName)
    // Writing .env file to current environment directory
    err := ioutil.WriteFile(filepath.Join(configDir, ".env"), []byte(content), 0777)
    return err 
}


