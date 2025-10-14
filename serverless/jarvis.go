package main

import (
	"jarvis/environment"
	"jarvis/internal/stacks"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/sirupsen/logrus"

	"github.com/aws/jsii-runtime-go"
)

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	environmentVars := &environment.Configuration{}

	err := environment.LoadEnvVars(environmentVars)

	if err != nil {
		logrus.Fatalf("Failed to load environment variables: %v", err)
	}

	logrus.Infof("Environment variables loaded successfully. %v", environmentVars)

	stacks.NewRegisterStack(app, "Register", &stacks.RegisterStackProps{
		StackProps: awscdk.StackProps{
			Description: jsii.String("Registers Jarvis's Application Commands with Discord"),
			Env:         env(),
		},
		AppID:  environmentVars.Discord.AppID,
		Token:  environmentVars.Discord.Token,
		ApiUrl: environmentVars.Discord.APIUrl,
	})

	stacks.NewAPIStack(app, "Jarvis-Commands-API", &stacks.APIStackProps{
		StackProps: awscdk.StackProps{
			Description: jsii.String("API for handling Jarvis's Application Commands"),
			Env:         env(),
		},
		PublicKey: environmentVars.Discord.PublicKey,
		Token:     environmentVars.Discord.Token,
		OpenAIKey: environmentVars.OpenAI.Key,
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
