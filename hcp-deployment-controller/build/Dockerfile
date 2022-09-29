FROM registry.access.redhat.com/ubi7/ubi-minimal:latest

ENV OPERATOR=/usr/local/bin/hcp-deployment-controller \
        USER_UID=root \
        USER_NAME=hcp-deployment-controller 

# install operator binary
COPY _output/bin/hcp-deployment-controller  ${OPERATOR}
COPY bin /usr/local/bin

RUN mkdir /usr/local/kubeconfig
RUN  /usr/local/bin/user_setup
RUN chmod +x /usr/local/bin/entrypoint

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}