.PHONY: help lint test test-vm schema deadcheck docs build clean

help:
	@grep -E '^[a-z-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN{FS=":.*?## "};{printf "  %-12s %s\n",$$1,$$2}'

build: ## stube bauen
	go build -trimpath -ldflags "-s -w" -o bin/stube ./cmd/stube

lint: ## Alles, was ohne Modell prüfbar ist
	golangci-lint run ./...
	yamllint hardware/ catalog/ .github/
	markdownlint docs/ *.md
	shellcheck $$(git ls-files '*.sh')
	vale docs/de/

test: ## Unit-Tests
	go test ./... -race -count=1

test-vm: ## Integrationstests gegen frische VMs. Langsam, nicht bei jedem Lauf.
	go test ./internal/... -tags=integration -timeout=45m

schema: ## Beispielmanifeste gegen das Schema prüfen
	go run ./cmd/stube validate --schema schema/homelab.schema.json testdata/manifests/*.yaml

deadcheck: ## Verwaiste Vorlagen, ungenutzte Felder, Doku ohne Code
	go run ./internal/tools/deadcheck

docs: ## Doku bauen, Verweise und Frontmatter prüfen
	go run ./internal/tools/docscheck
clean:
	rm -rf bin/ dist/
