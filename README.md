docker run -d \
  --name go-client \
  -p 4000:4000 \
  -e HOST_GRPC="IP:9000" \
  neftxx/go-client-grpc
