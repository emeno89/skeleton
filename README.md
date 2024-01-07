#### Shared folders

- local containers (docker-compose): [deploy/containers](deploy/containers);

- common kube resources: [deploy/kube](deploy/kube);

- common protobuf schemas: [proto/schema](proto/schema);

- common vendor: [vendor](vendor);

- common libraries: [shared](shared);

- common template .gitlab-ci: [.common-ci-template.yml](deploy/.common-ci-template.yml).

#### Services

- example GraphQL API service: [srv-gql](srv-gql);

- example GRPC service: [srv-grpc](srv-grpc);

#### Logging

Using [zap logger](https://github.com/uber-go/zap) with [ECS](https://go.elastic.co/ecszap) formatting
and wrapper [APM](https://go.elastic.co/apm/module/apmzap/v2). 
Logger object is created [here](shared/bootstrap/logger.go).

LOG_LEVEL и DEBUG env variables can be used to change log level (DEBUG=1 forcibly switches on debug level).

**NOTE:**
For local development you can use env template, for example: [srv-gql/cmd/.env.dst](srv-gql/cmd/.env.dst)

#### GraphQL schema generation:

```shell script
make gen_gql srv=<service_name>
````

#### Protobuf generation:

Install:

```shell script
sh proto/scripts/install.sh
````

Generate:

```shell script
make gen_pb
````

Hook for proto regeneration after git branch checkout:

```shell script
touch .git/hooks/post-checkout && \
chmod +x .git/hooks/post-checkout && \
echo "sh ./proto/scripts/generate.sh || true" > .git/hooks/post-checkout
````

#### Localization:

Install:

```shell script
go install github.com/nicksnyder/go-i18n/v2/goi18n@latest
````

1) generating texts from code;
2) applying translations.

```shell script
make gen_i18 srv=<service_name>
make apply_i18 srv=<service_name>
````

#### Linter:

Using [golangci-lint](https://golangci-lint.run/):

```shell script
make lint srv=<service_name>
````

#### Start docker-compose:

Preparations:

```shell script
make compose_build
make compose_prepare
````

**NOTE:**
If you use MacOS, ARCH can be changed to arm64 in [.env](deploy/containers/.env).

Start and get status:

```shell script
make compose_up
make compose_status
````