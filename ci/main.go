package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"dagger.io/dagger"
)

func checkEnvVariables(variables ...string) []string {
	var failedVariables []string

	for _, variable := range variables {
		if os.Getenv(variable) == "" {
			failedVariables = append(failedVariables, variable)
		}
	}

	return failedVariables
}

func main() {
	timeout := 10 * time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	requiredVariables := []string{
		"SEATTLE_UTILITIES_USERNAME",
		"SEATTLE_UTILITIES_PASSWORD",
	}

	failedVariables := checkEnvVariables(requiredVariables...)

	if len(failedVariables) > 0 {
		for _, variable := range failedVariables {
			fmt.Printf("Environment variable %s is not set\n", variable)
		}
		os.Exit(1)
	}

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	source := client.Container().
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
