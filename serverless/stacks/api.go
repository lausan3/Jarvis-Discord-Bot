package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/constructs-go/constructs/v10"
)

type APIStackProps struct {
	awscdk.StackProps
}

func NewAPIStack(scope constructs.Construct, id string, props *APIStackProps) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, &props.StackProps)

	// Define your API stack resources here
	return stack
}
