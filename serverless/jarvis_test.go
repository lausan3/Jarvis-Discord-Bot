package main

import (
	"jarvis/environment"
	"jarvis/utils/middleware"
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

func TestDiscordValidate(t *testing.T) {
	// Example values (replace with actual test values)
	body := `{\"app_permissions\":\"562949953601536\",\"application_id\":\"1291575538699472966\",\"attachment_size_limit\":524288000,\"authorizing_integration_owners\":{},\"entitlements\":[],\"id\":\"1425682843555856565\",\"token\":\"aW50ZXJhY3Rpb246MTQyNTY4Mjg0MzU1NTg1NjU2NTo4M0k3d1Vma3VpaVFwM0d5bkJWejc0OElOZzIxUGJWbEdFMTBSc3NINUFranZHam5hY0VwU1ViZXBmWkNaUjhrQnRPUG5YT0hGbUVyWUpOTWh1c1dZbVVJc2YwQVR6OHdnNWNKcDV3dUptRHNpdVZMTFlpUGthVmFjN21nMlNQUA\",\"type\":1,\"user\":{\"avatar\":\"c6a249645d46209f337279cd2ca998c7\",\"avatar_decoration_data\":null,\"bot\":true,\"clan\":null,\"collectibles\":null,\"discriminator\":\"0000\",\"display_name_styles\":null,\"global_name\":\"Discord\",\"id\":\"643945264868098049\",\"primary_guild\":null,\"public_flags\":1,\"system\":true,\"username\":\"discord\"},\"version\":1}"`
	signature := "e748aaef422c2313193d5ae5e1f60640a2663e0556d10f05d135efcf334848ddb5d04ef453f23612af02e53f18f39535014896b0deec167c919f29c86e96390d"
	timestamp := "1759979678"
	pubkey := "a6b962757aa712a666a8dee1a2d0364b1d2bf9ff221d0a32718f8fa88c2442c9"

	isValid := middleware.ValidateDiscordSecurityHeaders(body, signature, timestamp, pubkey)
	if !isValid {
		t.Errorf("Expected valid signature, got invalid")
	} else {
		logrus.Info("Signature validation passed")
	}
}
