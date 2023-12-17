The task scheduler will run a query into the database to get the jobs due at a specific minute. Then all the due jobs or tasks will be enqueued to a distributed message queue such as SQS or RabbitMQ. First-in-first-out (FIFO) queue would be the best.

We should have a primary-secondary configuration for the task scheduler to remove the single point of failure. If the primary server fails, secondary will take over.
