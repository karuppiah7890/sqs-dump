# sqs-dump

This is a tool to dump all the messages available in a SQS queue - usually a dead letter queue

## Running

Just build the tool using `make` and then run it

```bash
make
```

```bash
./sqs-dump
```

`sqs-dump` is just a simple tool and it keeps running and writes messages to `messages.json` file

### Setup

#### AWS credentials

Create an IAM user which has access to `sqs:GetQueueAttributes` and `sqs:ReceiveMessage`
