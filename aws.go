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

const ErrAlreadyExists string = "ResourceAlreadyExistsException"

type AWSConfig struct {
	Region          string
	AccessKeyId     string
	SecretAccessKey string
	LogGroup        string
	LogStream       string
	BatchSize       int
}

func getCreds(cfg AWSConfig) *credentials.Credentials {
	prv := credentials.StaticProvider{
		Value: credentials.Value{
			AccessKeyID:     cfg.AccessKeyId,
			SecretAccessKey: cfg.SecretAccessKey,
		},
	}
	return credentials.NewCredentials(&prv)
}

func getClient(cfg AWSConfig) *cloudwatchlogs.CloudWatchLogs {
	ses := session.Must(session.NewSession())
	stp := &aws.Config{
		Region:      &cfg.Region,
		Credentials: getCreds(cfg),
	}
	svc := cloudwatchlogs.New(ses, stp)
	return svc
}

func prepareLogGroup(c *cloudwatchlogs.CloudWatchLogs, cfg AWSConfig) error {
	grp := cloudwatchlogs.CreateLogGroupInput{
		LogGroupName: &cfg.LogGroup,
	}
	_, err := c.CreateLogGroup(&grp)
	if err != nil {
		if ae, ok := err.(awserr.Error); ok {
			if ae.Code() == ErrAlreadyExists {
				return nil
			}
		}
		return err
	}
	return nil
}

func prepareLogStream(c *cloudwatchlogs.CloudWatchLogs, cfg AWSConfig) error {
	stream := cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  &cfg.LogGroup,
		LogStreamName: &cfg.LogStream,
	}
	_, err := c.CreateLogStream(&stream)
	if err != nil {
		if ae, ok := err.(awserr.Error); ok {
			if ae.Code() == ErrAlreadyExists {
				return nil
			}
		}
		return err
	}
	return nil
}

func initialize(cfg AWSConfig) *cloudwatchlogs.CloudWatchLogs {
	c := getClient(cfg)
	var err error

	if err = prepareLogGroup(c, cfg); err != nil {
		panic(err)
	}
	if err = prepareLogStream(c, cfg); err != nil {
		panic(err)
	}

	return c
}

type awsLogsEmitter struct {
	config AWSConfig
	client *cloudwatchlogs.CloudWatchLogs
	batch  []*cloudwatchlogs.InputLogEvent
}

func NewAwsEmitter(cfg AWSConfig) *awsLogsEmitter {
	return &awsLogsEmitter{
		config: cfg,
		client: initialize(cfg),
		batch:  make([]*cloudwatchlogs.InputLogEvent, 0, cfg.BatchSize),
	}
}

func (x *awsLogsEmitter) emit(evs events) error {
	for _, e := range evs {
		if len(x.batch) == x.config.BatchSize {
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

func (x *awsLogsEmitter) flush() error {
	fmt.Println("flushing!")
	if len(x.batch) == 0 {
		fmt.Println("empty batch, bailing")
		return nil
	}

	batch := cloudwatchlogs.PutLogEventsInput{
		LogEvents:     x.batch,
		LogGroupName:  &x.config.LogGroup,
		LogStreamName: &x.config.LogStream,
	}
	out, err := x.client.PutLogEvents(&batch)
	if err != nil {
		panic(err)
		return err
	}
	fmt.Println(out)

	fmt.Println("flush success, nerfing batch")
	x.batch = make([]*cloudwatchlogs.InputLogEvent, 0, x.config.BatchSize)
	return nil
}
