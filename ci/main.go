package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"dagger.io/dagger"
)

func main() {
	timeout := 10 * time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Check that required environment variables are set
	var failedVariables []string
	requiredVariables := []string{
		"SEATTLE_UTILITIES_USERNAME",
		"SEATTLE_UTILITIES_PASSWORD",
	}
	for _, variable := range requiredVariables {
		if os.Getenv(variable) == "" {
			failedVariables = append(failedVariables, variable)
		}
	}
	if len(failedVariables) > 0 {
		for _, variable := range failedVariables {
			fmt.Fprintf(os.Stderr, "environment variable %s is not set\n", variable)
		}
		os.Exit(1)
	}

	// Connect to Docker and run tests
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	source := client.Container().
		// https://github.com/microsoft/playwright/tags
		// https://mcr.microsoft.com/en-us/product/playwright/tags
		// https://github.com/microsoft/playwright/releases
		From("mcr.microsoft.com/playwright:v1.41.0-jammy").
		WithExec([]string{"npm", "install", "-g", "npm@latest"}).
		WithDirectory("/src", client.Host().Directory("."), dagger.ContainerWithDirectoryOpts{
			Include: []string{
				"tests/",
				"playwright.config.ts",
				"package.json",
				"yarn.lock",
			},
		}).
		WithWorkdir("/src").WithExec([]string{
		"yarn", "install",
	})

	secret_user := client.SetSecret("utilities-username-secret", os.Getenv("SEATTLE_UTILITIES_USERNAME"))
	secret_pass := client.SetSecret("utilities-password-secret", os.Getenv("SEATTLE_UTILITIES_PASSWORD"))

	runner := source.WithEnvVariable("CACHEBUSTER", time.Now().String()).
		WithSecretVariable("SEATTLE_UTILITIES_USERNAME", secret_user).
		WithSecretVariable("SEATTLE_UTILITIES_PASSWORD", secret_pass)

	playwright := runner.WithExec([]string{"npx", "playwright", "test"})

	// Gather test results and export them to the host
	directories := map[string]string{
		"data":              "./data",
		"playwright-report": "./playwright-report",
		"test-results":      "./test-results",
	}

	for srcPath, destPath := range directories {
		out, err := playwright.Directory("/src/"+srcPath).Export(ctx, destPath)
		if err != nil {
			fmt.Printf("error exporting directory %s: %v\n", srcPath, err)
			panic(err)
		}
		fmt.Println(out)
	}
}
