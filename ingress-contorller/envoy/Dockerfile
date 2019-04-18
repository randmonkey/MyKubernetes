FROM envoyproxy/envoy:latest
COPY envoy.yaml /etc/envoy/envoy.yaml
ENTRYPOINT ["envoy","-c","/etc/envoy/envoy.yaml", "--service-cluster", "cluster0", "--service-node", "node0"]
