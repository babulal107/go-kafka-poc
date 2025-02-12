# go-kafka-poc
Learning Kafka integration with GoLang

## Apache Zookeeper & Kafka Basic Concept:
### Zookeeper 🦓
   Apache Zookeeper is a distributed coordination service used for managing configuration, leader election, and synchronization across distributed systems. 
   It ensures that multiple nodes in a distributed system can maintain a consistent state.

    - Zookeeper acts as a centralized service for managing and coordinating Kafka brokers. 
    - Since Kafka is a distributed system, Zookeeper helps keep everything synchronized and fault-tolerant.

### Apache Kafka 🦜
   Apache Kafka is a distributed event streaming platform used for high-throughput, fault-tolerant, real-time data streaming. 
   It is designed to handle large amounts of data in a publish-subscribe (pub-sub) model.

### Producer
   Produce a message to Kafka Message Queue

### Consumer
   Read(Consume) a message from Kafka Queue 

## Step to run on local
  Run Kafka Queue Container:
```sh
    make run-kafka
```

Down Kafka Queue Container:
```sh
    make down-kafka
```

Run Consumer:
```sh
    make run-consumer
```

Run Producer:
```sh
    make run-producer
```