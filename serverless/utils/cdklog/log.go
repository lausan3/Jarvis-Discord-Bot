package cdklog

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

// Helper function for logging to CDK console
func Log(scope constructs.Construct, message string) {
	awscdk.NewCfnOutput(scope, jsii.String("LogOutput"), &awscdk.CfnOutputProps{
		Value:       jsii.String(message),
		Description: jsii.String("Log message"),
	})
}
