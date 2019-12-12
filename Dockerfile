FROM registry.access.redhat.com/ubi8/ubi-minimal

ENV LANG=en_US.utf8 \
    GOPATH=/go \
    GOCACHE=/tmp/.cache \
    PATH=$PATH:$GOPATH/bin \
    OPERATOR_SDK_VERSION=0.13.0

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
    && mkdir -p /go/bin && mkdir -p /go/src \
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
RUN echo $GOPATH/bin/operator-sdk \
    && curl -L https://github.com/operator-framework/operator-sdk/releases/download/v$OPERATOR_SDK_VERSION/operator-sdk-v$OPERATOR_SDK_VERSION-x86_64-linux-gnu -o $GOPATH/bin/operator-sdk \
    && chmod +x $GOPATH/bin/operator-sdk  \
    && mkdir -p $GOPATH/src/$GO_PACKAGE_PATH/

# Stage 3 - Install Jenkins Operator

WORKDIR $GOPATH/src/$GO_PACKAGE_PATH
ADD ./ $GOPATH/src/$GO_PACKAGE_PATH

RUN make go-dependencies && go mod tidy && go mod vendor
ENV GO111MODULE=auto
RUN make build
RUN cp build/_output/bin/jenkins-operator /usr/local/bin/jenkins-operator

CMD [ "/usr/local/bin/jenkins-operator" ]
