# event-strangler

Idempotent Event Processing for Distributed Systems.

## Scenario

A system raises multiple events and publishes them to a topic in a single transaction, for example, a typical ordering
system could raise the following events `OrderCreated`, `OrderItemCreated` when an order is created, `OrderItemCreated`
is raised when items are appended to an order and `OrderItemDeleted` when items are removed from an order to an
`orders` topic.

For this scenario, we wish to export orders to a CSV file and publish them to a blob storage provider such as Amazon S3,
this process should be event based and run asynchronously. A solution is to employ serverless compute such as AWS
Lambda's and subscribe the lambda function to the `orders` topic. AWS provides subscription filters to exclude messages
from delivery to topic subscribers.

### Infrastructure

An example of the CDK infrastructure for this solution is demonstrated below. This CDK code provisions an S3 bucket with
the name `orders-bucket`, a lambda function `export-order-to-s3` that handles exporting and publishing the CSV file to
S3 and setting up the SNS subscription to invoke the `export-order-to-s3` lambda function with a filter policy only
allowing `OrderCreated`, `OrderItemCreated` and `OrderItemDeleted` events.

```ts
import * as lambda from '@aws-cdk/aws-lambda';
import * as s3 from '@aws-cdk/aws-s3';
import * as sns from '@aws-cdk/aws-sns';
import * as subscriptions from '@aws-cdk/aws-sns-subscriptions';

const ordersBucket = new s3.Bucket(this, 'OrdersBucket', {
    bucketName: 'orders-bucket',
    blockPublicAccess: s3.BlockPublicAccess.BLOCK_ALL,
});

const lambdaFunction = new lambda.Function(this, 'ExportOrderToS3', {
    functionName: 'export-order-to-s3',
    runtime: lambda.Runtime.NODEJS_14_X,
    handler: 'index.handler',
    code: Code.fromAsset(path.join(__dirname, 'src')),
});
ordersBucket.grantWrite(lambdaFunction);

const ordersTopic = new sns.Topic(this, 'OrdersTopic', {
    topicName: 'orders',
});
ordersTopic.addSubscription(new subscriptions.LambdaSubscription(lambdaFunction, {
    filterPolicy: {
        event: sns.SubscriptionFilter.stringFilter({
            allowlist: [
                'OrderCreated',
                'OrderItemCreated',
                'OrderItemDeleted',
            ],
        }),
    },
}));
```

### Concerns and Resolution

Whilst the granular events are more aligned with business actions, they also raise a problem with this specific
use-case; we only wish to process orders once per transaction, rather than processing the order for each event raised in
the transaction.

In the above scenario, when an order is created with two items, three events are raised by the orders service
(`OrderCreated`, `OrderItemCreated`, `OrderItemCreated`) resulting in the `export-order-to-s3` lambda function being
invoked three times. To resolve this, `event-strangler` utilises a store that can be queried to determine whether the
event we are processing has already been handled.
