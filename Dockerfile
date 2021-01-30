FROM docker.io/library/debian:10-slim

ADD hlybot /
ENTRYPOINT ["/hlybot"]
