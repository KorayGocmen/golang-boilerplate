package aws

import (
	"fmt"
	"strings"
	"sync"
	"time"

	awsgo "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	cwltypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/google/uuid"
	"github.com/koraygocmen/golang-boilerplate/internal/context"
)

type Message struct {
	Timestamp time.Time
	Content   []byte
}

type Stream struct {
	GroupName  string
	StreamName string
	Stop       chan bool

	cloudwatchLogsClient *cloudwatchlogs.Client
	queue                []Message
	lock                 *sync.RWMutex
}

func (s *Stream) Write(content []byte) (int, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.queue = append(s.queue, Message{
		Timestamp: time.Now().UTC(),
		Content:   content,
	})
	return len(s.queue), nil
}

func logStream(cloudwatchLogsClient *cloudwatchlogs.Client) func(ctx context.Ctx, groupName, streamNameSuffix string) (*Stream, error) {
	return func(ctx context.Ctx, groupName, streamNameSuffix string) (*Stream, error) {
		if cloudwatchLogsClient == nil {
			err := fmt.Errorf("log stream error: cloudwatch logs client is nil, maybe aws is not initialized?")
			return nil, err
		}

		if groupName == "" {
			err := fmt.Errorf("log stream error: group name is empty")
			return nil, err
		}

		stream := &Stream{
			GroupName:  groupName,
			StreamName: streamName(streamNameSuffix),
			Stop:       make(chan bool, 1),

			cloudwatchLogsClient: cloudwatchLogsClient,
			queue:                []Message{},
			lock:                 &sync.RWMutex{},
		}

		logStreamsOut, err := cloudwatchLogsClient.DescribeLogStreams(ctx, &cloudwatchlogs.DescribeLogStreamsInput{
			LogGroupName:        &stream.GroupName,
			LogStreamNamePrefix: &stream.StreamName,
		})
		if err != nil {
			return nil, err
		}

		// Create the log stream if it does not exist.
		if logStreamsOut == nil || len(logStreamsOut.LogStreams) == 0 {
			_, err = cloudwatchLogsClient.CreateLogStream(ctx, &cloudwatchlogs.CreateLogStreamInput{
				LogGroupName:  &stream.GroupName,
				LogStreamName: &stream.StreamName,
			})
			if err != nil {
				return nil, err
			}
		}

		go func() {
			queue := []cwltypes.InputLogEvent{}

			for {
				select {
				case <-stream.Stop:
					return
				default:
					{
						// Check if the stream queue has any messages.
						// if so, add them to the processing queue.
						stream.lock.RLock()
						if len(stream.queue) > 0 {
							for _, message := range stream.queue {
								messageContent := string(message.Content)
								queue = append(queue, cwltypes.InputLogEvent{
									Message:   &messageContent,
									Timestamp: awsgo.Int64(message.Timestamp.UnixNano() / int64(time.Millisecond)),
								})
							}

							// Dump the queue.
							stream.queue = []Message{}
						}
						stream.lock.RUnlock()

						// Process the queue.
						if len(queue) > 0 {
							_, err := stream.cloudwatchLogsClient.PutLogEvents(ctx, &cloudwatchlogs.PutLogEventsInput{
								LogEvents:     queue,
								LogGroupName:  &stream.GroupName,
								LogStreamName: &stream.StreamName,
							})
							if err != nil {
								continue
							}

							// Dump the process queue.
							queue = []cwltypes.InputLogEvent{}
						}
					}
				}

				time.Sleep(5 * time.Second)
			}
		}()

		return stream, nil
	}
}

func streamName(streamNameSuffix string) string {
	if streamNameSuffix == "" {
		// Remove the dashes from the uuid to make the stream
		// name consistent with the aws naming convention.
		streamNameSuffix = strings.ReplaceAll((uuid.New().String()), "-", "")
	}
	return fmt.Sprintf("logger/api/%s", streamNameSuffix)
}
