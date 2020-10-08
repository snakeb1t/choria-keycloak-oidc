FROM alpine:3.7

ENV SERVICE_HOME=/opt/palantir/services/{{Product}}
ENV CHORIA_KEYCLOAK_OIDC_INSTALL_CONFIG=$SERVICE_HOME/var/conf/install.yml
ENV CHORIA_KEYCLOAK_OIDC_RUNTIME_CONFIG=$SERVICE_HOME/var/conf/runtime.yml

RUN mkdir -p /opt/palantir/binaries && mkdir -p $SERVICE_HOME

# the templating is a gödel mechanism. "(InputDistArtifacts Product "sls")" returns a slice that contains the paths to
# the input dist artifacts for the dister with the name "sls" for the current product (returned by "Product"). The
# "(index ... 0)" returns the first path in the returned list of artifacts (which the SLS dister guarantees is the path
# to the SLS distribution artifact). See https://github.com/palantir/godel/wiki/Docker for details.
ADD {{index (InputDistArtifacts Product "sls") 0}} /opt/palantir/binaries

RUN ln -s /opt/palantir/binaries/{{Product}}-*/service $SERVICE_HOME/ && \
    cp -r /opt/palantir/binaries/{{Product}}-*/var $SERVICE_HOME/var

# Expose service port
EXPOSE 8443
# Expose management port
EXPOSE 8444

WORKDIR $SERVICE_HOME

ENTRYPOINT ["service/bin/linux-amd64/{{Product}}"]
CMD ["server"]
