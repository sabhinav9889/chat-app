package queue

import (
    "context"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/sqs"
)

// SQSQueueClient is an implementation of the QueueClient interface for SQS
type SQSQueueClient struct {
    client *sqs.SQS
    queueURL string
}

// NewSQSQueueClient creates a new SQSQueueClient instance
func NewSQSQueueClient(queueURL string) *SQSQueueClient {
    sess := session.Must(session.NewSession())
    sqsClient := sqs.New(sess)
    return &SQSQueueClient{
        client: sqsClient,
        queueURL: queueURL,
    }
}

// PublishMessage sends a message to the SQS queue
func (s *SQSQueueClient) PublishMessage(ctx context.Context, queue string, message []byte) error {
    _, err := s.client.SendMessageWithContext(ctx, &sqs.SendMessageInput{
        QueueUrl:    aws.String(queue),
        MessageBody: aws.String(string(message)),
    })
    return err
}

// ConsumeMessages processes messages from the SQS queue
func (s *SQSQueueClient) ConsumeMessages(ctx context.Context, queue string, handler func(message []byte) error) error {
    resp, err := s.client.ReceiveMessageWithContext(ctx, &sqs.ReceiveMessageInput{
        QueueUrl:            aws.String(queue),
        MaxNumberOfMessages: aws.Int64(10),
        WaitTimeSeconds:     aws.Int64(20),
    })
    if err != nil {
        return err
    }

    for _, msg := range resp.Messages {
        err := handler([]byte(*msg.Body))
        if err != nil {
            return err
        }
    }

    return nil
}
