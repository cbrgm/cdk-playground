package main

import (
	"github.com/aws/aws-cdk-go/awscdk"

	"github.com/aws/aws-cdk-go/awscdk/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/awslambda"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
)

const serviceName = "DynamoSample"

type SimpleApigatewayLambdaDynamoDBStackProps struct {
	awscdk.StackProps
}

func NewSimpleApigatewayLambdaDynamoDBStack(scope constructs.Construct, id string, props *SimpleApigatewayLambdaDynamoDBStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	table := awsdynamodb.NewTable(stack, jsii.String("Table"), &awsdynamodb.TableProps{
		TableName:   jsii.String(serviceName),
		BillingMode: awsdynamodb.BillingMode_PAY_PER_REQUEST,
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("pk"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		PointInTimeRecovery: jsii.Bool(false),
	})

	getter := awslambda.NewFunction(stack, jsii.String("Getter"), &awslambda.FunctionProps{
		FunctionName: jsii.String("get-message"),
		Code:         awslambda.Code_FromAsset(jsii.String("bin/getter-lambda"), nil),
		Runtime:      awslambda.Runtime_GO_1_X(),
		Handler:      jsii.String("main"),
		Timeout:      awscdk.Duration_Seconds(jsii.Number(30)),
	})
	table.GrantReadData(getter)

	setter := awslambda.NewFunction(stack, jsii.String("Setter"), &awslambda.FunctionProps{
		FunctionName: jsii.String("set-message"),
		Code:         awslambda.Code_FromAsset(jsii.String("bin/setter-lambda"), nil),
		Runtime:      awslambda.Runtime_GO_1_X(),
		Handler:      jsii.String("main"),
		Timeout:      awscdk.Duration_Seconds(jsii.Number(30)),
	})
	table.GrantReadWriteData(setter)

	createApiGateway(stack, getter, setter)

	return stack
}

func createApiGateway(stack awscdk.Stack, getter, setter awslambda.Function) {

	api := awsapigateway.NewRestApi(stack, jsii.String("ApiGateway"), &awsapigateway.RestApiProps{
		RestApiName: jsii.String(serviceName),
	})

	api.Root().AddResource(jsii.String("message/{proxy+}"), nil).AddMethod(
		jsii.String("GET"), awsapigateway.NewLambdaIntegration(getter, nil), nil)

	api.Root().AddResource(jsii.String("message/{proxy+}"), nil).AddMethod(
		jsii.String("POST"), awsapigateway.NewLambdaIntegration(setter, nil), nil)
}

func main() {
	app := awscdk.NewApp(nil)

	NewSimpleApigatewayLambdaDynamoDBStack(app, "SimpleApigatewayLambdaDynamoDBStack", &SimpleApigatewayLambdaDynamoDBStackProps{
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
