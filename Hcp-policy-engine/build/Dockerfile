FROM registry.access.redhat.com/ubi7/ubi-minimal:latest

ENV OPERATOR=/usr/local/bin/hcp-policy-engine \
    USER_UID=1001 \
    USER_NAME=hcp-policy-engine

# install operator binary
COPY _output/bin/hcp-policy-engine ${OPERATOR}

COPY bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]
CMD exec /bin/bash -c "trap : TERM INT; sleep infinity & wait"

USER ${USER_UID}
