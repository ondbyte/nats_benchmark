FROM scratch
COPY --from=nats:2.9.10-alpine3.17 /usr/local/bin/nats-server /nats-server
COPY test.config /test.config
EXPOSE 4222 8222 6222
ENV PATH="$PATH:/"
ENTRYPOINT ["/nats-server"]
CMD ["--js","--config", "test.config"]