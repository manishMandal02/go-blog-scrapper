.PHONY: tw-watch
tw-watch:
	./tailwindcss -i ./static/css/input.css -o ./static/css/style.css --watch

.PHONY: tailwind-build
tailwind-build:
	./tailwindcss -i ./static/css/input.css -o ./static/css/style.min.css --minify

.PHONY: templ-generate
templ-generate:
	templ generate

.PHONY: templ-watch
templ-watch:
	templ generate --watch  --proxy="http://localhost:8080" --open-browser=false
	
.PHONY: dev
dev:
	go build -o ./tmp/main ./cmd/scrapper/main.go && air

.PHONY: build
build:
	make tailwind-build
	make templ-generate
	go build -o ./bin/scrapper ./cmd/scrapper/main.go

.PHONY: vet
vet:
	go vet ./...

.PHONY: staticcheck
staticcheck:
	staticcheck ./...

.PHONY: test
test:
	  go test -race -v -timeout 30s ./..