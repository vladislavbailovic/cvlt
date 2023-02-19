package main

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

var (
	_awsRegion          string = "us-east-2"
	_awsAccessKeyId     string = "AKIA...."
	_awsSecretAccessKey string = "..."
	_awsLogGroup        string = "whatever"
	_awsLogStream       string = "test-test"
	_awsBatchSize       int    = 10
)

// TODO: creds auth
func getCreds() *credentials.Credentials {
	prv := credentials.StaticProvider{
		Value: credentials.Value{
			AccessKeyID:     _awsAccessKeyId,
			SecretAccessKey: _awsSecretAccessKey,
		},
	}
	return credentials.NewCredentials(&prv)
}

// TODO: region
func getClient() *cloudwatchlogs.CloudWatchLogs {
	ses := session.Must(session.NewSession())
	cfg := &aws.Config{
		Region:      &_awsRegion,
		Credentials: getCreds(),
	}
	svc := cloudwatchlogs.New(ses, cfg)
	return svc
}

// TODO: log group name
func prepareLogGroup(c *cloudwatchlogs.CloudWatchLogs) error {
	grp := cloudwatchlogs.CreateLogGroupInput{
		LogGroupName: &_awsLogGroup,
	}
	_, err := c.CreateLogGroup(&grp)
	if err != nil {
		if ae, ok := err.(awserr.Error); ok {
			if ae.Code() == "ResourceAlreadyExistsException" {
				return nil
			}
		}
		return err
	}
	return nil
}

// TODO: log group & stream name
func prepareLogStream(c *cloudwatchlogs.CloudWatchLogs) error {
	stream := cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  &_awsLogGroup,
		LogStreamName: &_awsLogStream,
	}
	_, err := c.CreateLogStream(&stream)
	if err != nil {
		if ae, ok := err.(awserr.Error); ok {
			if ae.Code() == "ResourceAlreadyExistsException" {
				return nil
			}
		}
		return err
	}
	return nil
}

func initialize() *cloudwatchlogs.CloudWatchLogs {
	c := getClient()
	var err error

	if err = prepareLogGroup(c); err != nil {
		panic(err)
	}
	if err = prepareLogStream(c); err != nil {
		panic(err)
	}

	return c
}

type awsLogsEmitter struct {
	client *cloudwatchlogs.CloudWatchLogs
	batch  []*cloudwatchlogs.InputLogEvent
}

// TODO batch queue size
func NewAwsEmitter() *awsLogsEmitter {
	return &awsLogsEmitter{
		client: initialize(),
		batch:  make([]*cloudwatchlogs.InputLogEvent, 0, _awsBatchSize),
	}
}

func (x *awsLogsEmitter) emit(evs events) error {
	for _, e := range evs {
		if len(x.batch) == _awsBatchSize {
			fmt.Println("update: triggering flush")
			x.flush()
		}
		now := time.Now().UnixMicro() / 1000
		msg := e.Entry()
		logEvent := cloudwatchlogs.InputLogEvent{
			Timestamp: &now,
			Message:   &msg,
		}
		fmt.Printf("\t- creating new event: %+v\n", logEvent)
		x.batch = append(x.batch, &logEvent)
		fmt.Println("new batch size", len(x.batch))
	}
	return nil
}

// TODO _aws*
func (x *awsLogsEmitter) flush() error {
	fmt.Println("flushing!")
	if len(x.batch) == 0 {
		fmt.Println("empty batch, bailing")
		return nil
	}

	batch := cloudwatchlogs.PutLogEventsInput{
		LogEvents:     x.batch,
		LogGroupName:  &_awsLogGroup,
		LogStreamName: &_awsLogStream,
	}
	out, err := x.client.PutLogEvents(&batch)
	if err != nil {
		panic(err)
		return err
	}
	fmt.Println(out)

	fmt.Println("flush success, nerfing batch")
	x.batch = make([]*cloudwatchlogs.InputLogEvent, 0, _awsBatchSize)
	return nil
}
