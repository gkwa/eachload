package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"dagger.io/dagger"
)

func main() {
	timeout := 3 * time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if os.Getenv("SEATTLE_UTILITIES_USERNAME") == "" {
		panic("Environment variable SEATTLE_UTILITIES_USERNAME is not set")
	}
	if os.Getenv("SEATTLE_UTILITIES_PASSWORD") == "" {
		panic("Environment variable SEATTLE_UTILITIES_PASSWORD is not set")
	}

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	runner := client.Container().
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

	// I want playwrite test to run everytime without caching.  https://docs.dagger.io/cookbook/#invalidate-cache
	out, err := runner.WithEnvVariable("CACHEBUSTER", time.Now().String()).
		WithSecretVariable("SEATTLE_UTILITIES_USERNAME", secret_user).
		WithSecretVariable("SEATTLE_UTILITIES_PASSWORD", secret_pass).
		WithExec([]string{"npx", "playwright", "test"}).
		Directory("/src/data").
		Export(ctx, "./data")
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
}
