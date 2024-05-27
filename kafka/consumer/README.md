# ca-go/kafka/consumer

intro

## Environment Variables

You can optionally set these environment variables instead of passing these are options to `NewConsumer()`
To use the package level methods `Encode` and `Decode` you MUST set these:

- AUTH_PUBLIC_JWK_KEYS = A JSON string containing the well known public keys for Decoding a token.

- KAFKA_BROKERS = A list of Kafka brokers. Calling `WithBrokers(brokers)` overwrites this default.
- KAFKA_TOPICS = A list of Kafka topics to consumer. Calling `WithTopics(topics)` overwrites this default.

## Service

## Consumer

## Examples
