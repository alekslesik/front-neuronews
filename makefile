# Change these variables as necessary.
MAIN_PATH := ./cmd/neuro-news
BINARY_NAME := neuro-news

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

## tar: archieve project
.PHONY: tar
tar:
	tar -czpf ${BINARY_NAME}.tar.gz *

## untar: retreive project
.PHONY: untar
untar:
	tar -xzpf ${BINARY_NAME}.tar.gz -C .

## go111on: set GO111MODULE='on'
.PHONY: go111on
go111on:
	go env -w GO111MODULE='on'

## go111off: set GO111MODULE='off'
.PHONY: go111off
go111off:
	go env -w GO111MODULE='off'

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run: run app
.PHONY: run
run:
	go run $(MAIN_PATH)/main.go

## run-build: run and build app local
.PHONY: run-build
run-build: build
	./build/$(BINARY_NAME)

## build: build app local
.PHONY: build
build: tidy
	go build -o ./build $(MAIN_PATH)/

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

## audit: run quality control checks
.PHONY: audit
audit: tidy
	go mod verify
	go vet ./...

	# Линтер для исходного кода на Go. Он предлагает советы по стилю кодирования,
	# соответствующие Руководству по стилю кода Go.
	# golint ./...

	# Инструмент для анализа безопасности кода на Go.Он проверяет код на наличие потенциальных уязвимостей,
	# таких как SQL-инъекции, небезопасное использование криптографии
	# gosec ./...

	# Инструмент для обнаружения мест, где ошибки, возвращаемые функциями, не обрабатываются.
	# В Go принято обрабатывать почти все ошибки,
	# и игнорирование ошибок может привести к непредсказуемому поведению программы
	# errcheck ./...

	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go test -race -buildvcs -vet=off ./...

# ==================================================================================== #
# TESTING
# ==================================================================================== #
## test-cover: test cover of whole project
.PHONY: test-cover
test-cover:
	go test -v -race -coverprofile=./tests/coverage.out ./...
	go tool cover -html=./tests/coverage.out -o ./tests/coverage.html

# ==================================================================================== #
# MYSQL
# ==================================================================================== #

## mysql-root: connect to mySQL by root user
.PHONY: mysql-root
mysql-root:
	mysql -u root -p

# ==================================================================================== #
# DEPLOY
# ==================================================================================== #

## deploy: deploy to host by ansible
.PHONY: deploy
deploy: build
	ansible-playbook -i ops/production/ansible/hosts.ini ops/production/ansible/dpl.yml -vv

# ==================================================================================== #
# DOCKER
# ==================================================================================== #

## docker-run: run docker container with app
.PHONY: docker-run
docker-run: build
	ansible-playbook -i ops/production/ansible/hosts.ini ops/production/ansible/dpl.yml -vv

# ==================================================================================== #
# OTHER
# ==================================================================================== #

## go-update: update Golang, use "v" for version (go-update v=1.21.5)
.PHONY: go-update
go-update:
	sudo rm -rf go$(v).linux-amd64.tar.gz
	wget https://go.dev/dl/go$(v).linux-amd64.tar.gz
	sudo rm -rf /usr/local/go
	sudo tar -C /usr/local -xzf go$(v).linux-amd64.tar.gz
	sudo rm -rf go$(v).linux-amd64.tar.gz
	go version