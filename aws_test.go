package main

import (
	"strings"
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

func Test_EventConversion(t *testing.T) {
	suite := map[string]string{
		"2023-01-08T02:00:00Z": "old",
		"2032-01-08T02:00:00Z": "future",
	}
	for ts, want := range suite {
		t.Run(ts, func(t *testing.T) {
			event := jsonLogEvent{Time: ts, Log: "wat"}
			_, err := logEvent2InputLogEvent(event)
			if err == nil {
				t.Error("expected error")
			}
			if !strings.Contains(err.Error(), want) {
				t.Errorf("want %q, got %q",
					want, err)
			}
		})
	}
}
