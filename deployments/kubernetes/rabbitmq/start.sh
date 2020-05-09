#!/bin/bash

# Install RabbitMQ
kubectl apply -f ./secret.yml
helm install rabbitmq bitnami/rabbitmq -f ./values.yml --set rabbitmq.username=guest,rabbitmq.password=guest