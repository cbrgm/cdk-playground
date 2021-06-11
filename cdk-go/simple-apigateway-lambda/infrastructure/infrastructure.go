package main

import (
	"github.com/aws/aws-cdk-go/awscdk"

	"github.com/aws/aws-cdk-go/awscdk/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/awslambda"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
)

const serviceName = "helloWorld"

type SimpleApigatewayLambdaStackProps struct {
	awscdk.StackProps
}

func NewSimpleApigatewayLambdaStack(scope constructs.Construct, id string, props *SimpleApigatewayLambdaStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	helloWorld := awslambda.NewFunction(stack, jsii.String("HelloWorld-Handler"), &awslambda.FunctionProps{
		FunctionName: jsii.String("hello-world-handler"),
		Code:         awslambda.Code_FromAsset(jsii.String("bin/hello-world"), nil),
		Runtime:      awslambda.Runtime_GO_1_X(),
		Handler:      jsii.String("main"),
		Timeout:      awscdk.Duration_Seconds(jsii.Number(30)),
		Environment: &map[string]*string{
			"HELLO_MESSAGE": jsii.String("Hello World!"),
		},
	})

	createApiGateway(stack, helloWorld)

	return stack
}

func createApiGateway(stack awscdk.Stack, helloWorldLambda awslambda.Function) {

	api := awsapigateway.NewRestApi(stack, jsii.String("ApiGateway"), &awsapigateway.RestApiProps{
		RestApiName: jsii.String(serviceName),
	})

	api.Root().AddResource(jsii.String("hello"), nil).AddMethod(
		jsii.String("GET"), awsapigateway.NewLambdaIntegration(helloWorldLambda, nil), nil)
}

func main() {
	app := awscdk.NewApp(nil)

	NewSimpleApigatewayLambdaStack(app, "SimpleApigatewayLambdaStack", &SimpleApigatewayLambdaStackProps{
		awscdk.StackProps{
			Env: env(),
		},
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
