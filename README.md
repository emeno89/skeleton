#### Служебные разделы

- Контейнеры для запуска локально (docker-compose): [deploy/containers](deploy/containers);

- Общие ресурсы для работы kubernetes: [deploy/kube](deploy/kube);

- Все файлы схемы взаимодействия proto grpc: [proto/schema](proto/schema);

- Общий вендор для всех микросервисов: [vendor](vendor);

- Общий код для всех микросервисов не из вендора: [shared](shared);

- Общий шаблон .gitlab-ci: [.common-ci-template.yml](deploy/.common-ci-template.yml).

#### Сервисы

- Логика GraphQL API: [srv-gql](srv-gql);

- Логика GRPC-сервиса: [srv-grpc](srv-grpc);

#### Логгирование

Используется [zap logger](https://github.com/uber-go/zap) с добавлением
форматирования [ECS](https://go.elastic.co/ecszap)
и трейсинга [APM](https://go.elastic.co/apm/module/apmzap/v2). Настройка происходит [здесь](shared/bootstrap/logger.go).

Для управления уровнем логгирования используется env-переменная LOG_LEVEL и DEBUG (для принудительного включения debug).

**NOTE:**
Для локального запуска можно использовать шаблон env-переменных, к примеру: [srv-gql/cmd/.env.dst](srv-gql/cmd/.env.dst)

#### Генерация из схемы GraphQL:

```shell script
make gen_gql srv=<service_name>
````

#### Генерация proto GRPC:

Устанавливаем:

```shell script
sh proto/scripts/install.sh
````

Генерация proto:

```shell script
make gen_pb
````

Хук на перегенерацию proto при переключении между ветками:

```shell script
touch .git/hooks/post-checkout && \
chmod +x .git/hooks/post-checkout && \
echo "sh ./proto/scripts/generate.sh || true" > .git/hooks/post-checkout
````

#### Локализация:

Устанавливаем:

```shell script
go install github.com/nicksnyder/go-i18n/v2/goi18n@latest
````

#### Линтер:

Используем [golangci-lint](https://golangci-lint.run/):

```shell script
make lint srv=<service_name>
````

1) генерация текстов из кода;
2) применение переводов.

```shell script
make gen_i18 srv=<service_name>
make apply_i18 srv=<service_name>
````

#### Запуск docker-compose:

Выполняем подготовительный запуск:

```shell script
make compose_build
make compose_prepare
````

**NOTE:**
Если запуск происходит на MacOS, возможно требуется поменять в [.env](deploy/containers/.env) ARCH на arm64.

Запускаем и проверяем:

```shell script
make compose_up
make compose_status
````