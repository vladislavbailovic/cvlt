package main

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws/credentials"
)

func getTestFailCfg() AWSConfig {
	return AWSConfig{
		Region:          "us-east-2",
		AccessKeyId:     "AKIA....",
		SecretAccessKey: "...",
		LogGroup:        "whatever",
		LogStream:       "test-test",
		BatchSize:       10,
	}
}

func Test_getCreds(t *testing.T) {
	creds := getCreds(getTestFailCfg())
	val, err := creds.Get()
	if err != nil {
		t.Error(err)
	}
	credsTest(val, t)
}

func credsTest(val credentials.Value, t *testing.T) {
	if val.ProviderName != "StaticProvider" {
		t.Errorf("want static provider, got %q",
			val.ProviderName)
	}
	if val.AccessKeyID == "" {
		t.Error("wanted access key ID")
	}
	if val.SecretAccessKey == "" {
		t.Error("wanted secret access key")
	}
}

func Test_getClient(t *testing.T) {
	c := getClient(getTestFailCfg())
	val, err := c.Config.Credentials.Get()
	if err != nil {
		t.Error(err)
	}
	credsTest(val, t)
}

// func Test_prepareGroupAndStream_HappyPath(t *testing.T) {
// 	c := initialize(getTestFailCfg())
// 	val, err := c.Config.Credentials.Get()
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	credsTest(val, t)
// }
