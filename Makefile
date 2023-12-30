gen_pb:
	sh ./proto/scripts/generate.sh
gen_gql:
	cd ${srv}/internal/api/graphql/graph && go run github.com/99designs/gqlgen generate
gen_i18:
	cd ${srv} && goi18n extract -sourceLanguage ru -outdir ./cmd/translation -format yaml
	cd ${srv} && goi18n merge -sourceLanguage ru -outdir ./cmd/translation -format yaml ./cmd/translation/active.*.yaml
apply_i18:
	cd ${srv} && goi18n merge -sourceLanguage ru -outdir ./cmd/translation -format yaml ./cmd/translation/active.*.yaml ./cmd/translation/translate.*.yaml
compose_prepare:
	cd deploy/containers && cp -r .app_env.dst .app_env && cp .env.dst .env && mkdir -p volumes
compose_build:
	cd deploy/containers && docker-compose build $(srv)
compose_up:
	cd deploy/containers && docker-compose up -d $(srv) --remove-orphans
compose_down:
	cd deploy/containers && docker-compose down --remove-orphans
compose_status:
	cd deploy/containers && docker-compose ps
test:
	go test ./$(srv)/... -mod=vendor -v -timeout 10s | grep -v "no test files"
lint:
	golangci-lint run -v ./$(srv)/...