package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2integrations"
	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type APIStackProps struct {
	awscdk.StackProps

	PublicKey string
	Token     string

	OpenAIKey string
}

func NewAPIStack(scope constructs.Construct, id string, props *APIStackProps) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, &props.StackProps)
	publicKey := props.PublicKey
	token := props.Token
	oaiKey := props.OpenAIKey

	api := awsapigatewayv2.NewHttpApi(stack, jsii.String("JarvisAPI"), &awsapigatewayv2.HttpApiProps{
		ApiName:     jsii.String("JarvisBot-CommandsAPI"),
		Description: jsii.String("API for handling Jarvis's Application Commands"),
		CorsPreflight: &awsapigatewayv2.CorsPreflightOptions{
			AllowHeaders: jsii.Strings("Content-Type", "Authorization", "X-Signature-Ed25519", "X-Signature-Timestamp"),
			AllowMethods: &[]awsapigatewayv2.CorsHttpMethod{
				awsapigatewayv2.CorsHttpMethod_POST,
				awsapigatewayv2.CorsHttpMethod_OPTIONS,
			},
			AllowOrigins: jsii.Strings("*"),
		},
	})

	commandFunc := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("JarvisCommandsHandler"), &awscdklambdagoalpha.GoFunctionProps{
		Entry:      jsii.String("internal/lambdas/api/commands.go"),
		MemorySize: jsii.Number(128),
		Timeout:    awscdk.Duration_Seconds(jsii.Number(10)),
		Environment: &map[string]*string{
			"DISCORD_PUBLIC_KEY": &publicKey,
			"DISCORD_BOT_TOKEN":  &token,
			"OPENAI_API_KEY":     &oaiKey,
		},
	})

	commandIntegration := awsapigatewayv2integrations.NewHttpLambdaIntegration(
		jsii.String("JarvisCommandsIntegration"),
		commandFunc,
		&awsapigatewayv2integrations.HttpLambdaIntegrationProps{},
	)

	api.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Path: jsii.String("/interactions"),
		Methods: &[]awsapigatewayv2.HttpMethod{
			awsapigatewayv2.HttpMethod_POST,
		},
		Integration: commandIntegration,
	})

	return stack
}
