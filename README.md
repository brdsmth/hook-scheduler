## Introduction

The `hook-scheduler` service primarily runs queries on the instance of **DynamoDB** to find jobs that need to be executed. This is accomplished by indexing the jobs on in the `queue` table and comparing their `ExecuteAt` timestamp to the current timestamp. The `hook-scheduler` performs these queries by method of a **cronjob** that runs every minute. If the `ExecuteAt` timestamp is earlier than the current time then the `hook-scheduler` will take that job and add it to the **RabbitMQ** distributed message queue for processing.

The use of a cronjob running every minute ensures that we can provide accuracy of the `ExecuteAt` time with a margin of error of one minute.

We should have a primary-secondary configuration for the task scheduler to remove the single point of failure. If the primary server fails, secondary will take over.

### Usage

The `hook-scheduler` service is meant to be run in conjunction with the `hook-api` and the `hook-runner` services.

To run the `hook-scheduler` service the following `.env` variables need to be set

```
AWS_CONFIG_PROFILE=
RABBITMQ_URL=
RABBITMQ_QUEUE=
DYNAMODB_QUEUE_TABLE=
```

Once these are added the server can be started by

```
go run main.go
```

Once the `hook-scheduler` server is running the cron service will start. Currently, the cron service is set to run every minute(`* * * * *`).
