### consul && consul-template ###
FROM hashicorp/consul-template:latest AS consul-template
FROM consul:latest AS consul

### build app and differ binaries ###
FROM golang:latest as go
RUN go get github.com/spf13/viper
COPY app /app
WORKDIR /app
RUN go build -o app main.go
COPY differ /differ
WORKDIR /differ
RUN go build -o differ main.go

### demo image ###
FROM ansible/centos7-ansible:latest
COPY demo_container_src /demo
RUN mv /demo/runit.py /opt/ansible/ansible/lib/ansible/modules/system/runit.py

RUN rpm -i /demo/runit-2.1.1-7.el7.centos.x86_64.rpm && \
      rm /demo/runit-2.1.1-7.el7.centos.x86_64.rpm

WORKDIR /demo/app
COPY --from=go /app/app bin/app
COPY app/config.yaml .
WORKDIR /demo/differ
COPY --from=go /differ/differ bin/differ

COPY --from=consul-template /consul-template /demo/consul-template/bin/
COPY --from=consul /bin/consul /bin/
RUN ln -s /demo/consul-template/runit_service /etc/service/consul-template

WORKDIR /demo
ENTRYPOINT ["/sbin/runsvdir-start"]
