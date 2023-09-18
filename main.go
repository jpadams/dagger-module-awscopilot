package main

import (
	"context"
	"fmt"
	"runtime"
)

type Awscopilot struct{}

func (m *Awscopilot) Base(ctx context.Context, alpineVersion string, copilotVersion string) (*Container, error) {
	os := runtime.GOOS 
	if os != "linux" {
		panic("Only Linux containers are supported!")
	}
	arch := runtime.GOARCH
	archSuffix := ""
	if arch == "arm64" {
		archSuffix = "-arm64"
	} else if arch != "amd64" {
		panic("Only arm64 and amd64 architectures are supported!")
	}
	installcmd := fmt.Sprintf("curl -s -Lo copilot https://github.com/aws/copilot-cli/releases/%s/download/copilot-%s%s && chmod +x copilot && mv copilot /usr/local/bin/copilot", copilotVersion, os, archSuffix)
	alpine := fmt.Sprintf("alpine:%s", alpineVersion)
	//image := fmt.Sprintf("public.ecr.aws/docker/library/docker:%s", alpineVersion)
	//dind := dag.Container().From(image). 	
	//	WithExec([]string{"dockerd", "--host=tcp://0.0.0.0:2375", "--tls=false"}, ContainerWithExecOpts{InsecureRootCapabilities: true, ExperimentalPrivilegedNesting: true}).
	//	WithExposedPort(2375)
	return dag.Container().
		From(alpine).
		WithExec([]string{"apk", "add", "curl", "git"}).
		WithExec([]string{"sh", "-c", installcmd}).
	//	WithServiceBinding("dind", dind).
	//	WithEnvVariable("DOCKER_HOST", "dind:2375").Sync(ctx)
		Sync(ctx)
}

func (m *Awscopilot) Test(ctx context.Context, awsid string, awssecret string, awsregion string) (string, error) {
	c, err := (&Awscopilot{}).Base(ctx, "latest", "latest")
	if err != nil {
		panic(err)	
	}
	return c.WithExec([]string{"sh", "-c", "git clone https://github.com/aws-samples/aws-copilot-sample-service example"}).
		WithWorkdir("example").
		WithSecretVariable("AWS_ACCESS_KEY_ID", dag.SetSecret("awsid", awsid)).
		WithSecretVariable("AWS_SECRET_ACCESS_KEY", dag.SetSecret("awssecret", awssecret)).
		WithEnvVariable("AWS_REGION", awsregion).
		//WithExec([]string{"sh", "-c", "docker ps"}).
		//WithExec([]string{"sh", "-c", "copilot env init --name test --profile default --default-config"}).
//		WithExec([]string{"sh", "-c", "copilot env init --name test --profile default --default-config"}).
		//WithExec([]string{"sh", "-c", "copilot init --app demo --name api --type \"Load Balanced Web Service\" --dockerfile \"./Dockerfile\" --deploy"}).
		//WithExec([]string{"sh", "-c", "copilot app delete --name api --yes"}).
//		WithExec([]string{"sh", "-c", "copilot env delete --name test --yes"}).
		WithExec([]string{"sh", "-c", "copilot help"}).
		Stdout(ctx)
}
