package main

import (
	"jarvis/environment"
	"testing"

	"github.com/sirupsen/logrus"
)

// example tests. To run these tests, uncomment this file along with the
// example resource in serverless_test.go
// func TestJarvisStack(t *testing.T) {
// 	// GIVEN
// 	app := awscdk.NewApp(nil)

// 	// WHEN
// 	stack := NewJarvisStack(app, "MyStack", nil)

// 	// THEN
// 	template := assertions.Template_FromStack(stack, nil)

// 	template.HasResourceProperties(jsii.String("AWS::SQS::Queue"), map[string]interface{}{
// 		"VisibilityTimeout": 300,
// 	})
// }

func TestEnvVarLoading(t *testing.T) {
	config := &environment.Configuration{}

	err := environment.LoadEnvVars(config)
	if err != nil {
		t.Errorf("Failed to load environment variables: %v", err)
	}

	logrus.Infof("Loaded configuration: %v", config)
	logrus.Infof("Discord App ID: %s", config.Discord.AppID)
}
