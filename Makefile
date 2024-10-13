
all: pack-js
all: build-ws
all: copy-js
all: run-ws

node_modules:
	pnpm install 

pack-js: node_modules
	pnpm build

copy-js:
	cat build/bundle.js | pbcopy

build-ws: cert
	cp build/bundle.js ./cmd/server/sessionmanager/bundle.js
	go build ./cmd/server

run-ws:
	./server

cert:
	go test -v ./...

pack-codes:
	zip -r source.zip ./cmd/ ./src/ ./cert/ .gitignore go.mod go.sum Makefile package.json pnpm-lock.yaml README.md tsconfig.json webpack.config.js