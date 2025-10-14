package stacks

import (
	"time"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/customresources"
	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type RegisterStackProps struct {
	awscdk.StackProps

	AppID  string
	Token  string
	ApiUrl string
}

type RegisterStack struct {
	awscdk.Stack
}

// NewRegisterStack functionally acts as a script to register Discord commands using a Lambda function.
func NewRegisterStack(scope constructs.Construct, id string, props *RegisterStackProps) *RegisterStack {
	sprops := props.StackProps
	stack := awscdk.NewStack(scope, &id, &sprops)

	appID := props.AppID
	token := props.Token
	apiUrl := props.ApiUrl

	function := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("RegisterCommandsFunction"), &awscdklambdagoalpha.GoFunctionProps{
		Description: jsii.String("Lambda function to register Discord commands for Jarvis"),
		Entry:       jsii.String("internal/lambdas/register/register.go"),
		Environment: &map[string]*string{
			"DISCORD_APP_ID":    &appID,
			"DISCORD_BOT_TOKEN": &token,
			"DISCORD_API_URL":   &apiUrl,
		},
		Timeout: awscdk.Duration_Millis(jsii.Number(30 * time.Second.Milliseconds())), // 30 seconds
	})

	function.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("cloudformation:DeleteStack", "cloudformation:DescribeStacks"),
		Resources: jsii.Strings("*"),
	}))

	provider := customresources.NewProvider(stack, jsii.String("RegisterCommandsProvider"), &customresources.ProviderProps{
		OnEventHandler: function,
	})

	awscdk.NewCustomResource(stack, jsii.String("RegisterCommandsResource"), &awscdk.CustomResourceProps{
		ServiceToken: provider.ServiceToken(),
	})

	return &RegisterStack{
		Stack: stack,
	}
}
