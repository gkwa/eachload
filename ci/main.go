package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

func main() {
	ctx := context.Background()

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

	secret_user := client.SetSecret("utilities-username-secret", os.Getenv("SEATTLE_UTILITIES_USERNAME"))
	secret_pass := client.SetSecret("utilities-password-secret", os.Getenv("SEATTLE_UTILITIES_PASSWORD"))

	source := client.Container().
		From("mcr.microsoft.com/playwright:v1.41.0-jammy").
		WithSecretVariable("SEATTLE_UTILITIES_USERNAME", secret_user).
		WithSecretVariable("SEATTLE_UTILITIES_PASSWORD", secret_pass).
		WithDirectory("/src", client.Host().Directory("."), dagger.ContainerWithDirectoryOpts{
			Include: []string{
				"tests/", "playwright.config.ts",
				"package.json", "yarn.lock",
			},
		})

	runner := source.WithWorkdir("/src").WithExec([]string{
		"yarn", "install",
	})

	out, err := runner.WithExec([]string{"ls"}).Stderr(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(out)

	out, err = runner.WithExec([]string{"npx", "playwright", "test"}).Stderr(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
}
