docker exec -it kafka \
  kafka-topics --create \
    --topic studio-stream \
    --bootstrap-server localhost:9092 \
    --partitions 2 \
    --replication-factor 1
