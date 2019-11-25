FROM registry.access.redhat.com/ubi8/ubi-minimal

ENV LANG=en_US.utf8 \
    GOPATH=/tmp/go \
    GOCACHE=/tmp/.cache \
    PATH=$PATH:$GOPATH/bin

ARG GO_PACKAGE_PATH=github.com/redhat-developer/jenkins-operator

# Stage 1 - Install dependencies

RUN microdnf install -y \
    findutils \
    git \
    make \
    procps-ng \
    tar \
    wget \
    which \
    bc \
    gcc

WORKDIR /tmp

# download, verify and install golang
ENV PATH=$PATH:/usr/local/go/bin
RUN curl -Lo go1.13.4.linux-amd64.tar.gz https://dl.google.com/go/go1.13.4.linux-amd64.tar.gz \
    && echo "692d17071736f74be04a72a06dab9cac1cd759377bd85316e52b2227604c004c go1.13.4.linux-amd64.tar.gz" > go1.13.4.linux-amd64.sha256 \
    && sha256sum -c go1.13.4.linux-amd64.sha256 \
    && tar xzf go1.13.4.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go1.13.4.linux-amd64.tar.gz \
    && rm -f go1.13.4.linux-amd64.tar.gz \
    && go version

# download and install kubectl
RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl \
    && chmod +x ./kubectl \
    && mv ./kubectl /usr/local/bin/kubectl

# Stage 2 - Install operator-sdk
RUN echo $GOPATH/bin/operator-sdk
RUN curl -L https://github.com/operator-framework/operator-sdk/releases/download/v$OPERATOR_SDK_VERSION/operator-sdk-v$OPERATOR_SDK_VERSION-x86_64-linux-gnu -o $GOPATH/bin/operator-sdk \
    && chmod +x $GOPATH/bin/operator-sdk

RUN mkdir -p ${GOPATH}/src/${GO_PACKAGE_PATH}/

WORKDIR ${GOPATH}/src/${GO_PACKAGE_PATH}

ENTRYPOINT [ "/bin/bash" ]