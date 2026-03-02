.PHONY: run
run: ## Run run app.
	$(DOCKER_RUN) $(ROOT_ARGS)

.PHONY: run-help
run-help: ## Run help app.
	$(DOCKER_RUN) help

.PHONY: run-version
run-version: ## Run version app.
	$(DOCKER_RUN) version

.PHONY: run-conf
run-conf: ## Run configure app.
	$(DOCKER_RUN) configure

.PHONY: run-add
run-add: ## Run add app.
	$(DOCKER_RUN) add $(ADD_ARGS)

.PHONY: run-add-ux
run-add-ux: ## Run add app in an interactive mode.
	$(DOCKER_RUN) add -u -n default

.PHONY: run-add-ux-with-ns
run-add-ux-with-ns: ## Run add app in an interactive mode with predefined namespace.
	$(DOCKER_RUN) add -u -n $(KBNAMESPACE)

.PHONY: run-import
run-import: ## Run import app .
	$(DOCKER_RUN) import

.PHONY: run-import-sample
run-import-sample: ## Run import app to load sample kbs.
	$(DOCKER_RUN) import -f ../../docs/samples/import-sample.yaml

.PHONY: run-get
run-get: ## Run get app.
	$(DOCKER_RUN) get $(GET_ARGS)

.PHONY: run-get-ux
run-get-ux: ## Run get app with ux.
	$(DOCKER_RUN) get -u

.PHONY: run-update
run-update: ## Run update kb app.
	$(DOCKER_RUN) update

.PHONY: run-add-with-args
run-add-with-args: ## Run add app with predefined arguments.
	$(DOCKER_RUN) add \
	-k btc -v crypto -o currencies -c crypto \
	-t btc,crypto,currencies,blockchain \
	-r dementor -n personal

.PHONY: run-add-media
run-add-media: ## Run add app to save kb with media.
	$(DOCKER_RUN) add \
	-k btc -o currencies -c media \
	-t btc,crypto,currencies,blockchain \
	-r dementor -n personal \
	-v 'https://pbs.twimg.com/media/GZL9kSeXgAAXf3B?format=jpg&name=4096x4096'

.PHONY: run-import-with-args
run-import-with-args: ## Run import app with predefined arguments.
	$(DOCKER_RUN) import -f ../../docs/samples/import-sample.yaml --show-added-kbs --show-failed-kbs

.PHONY: run-export-with-args
run-export-with-args: ## Run export app with predefined arguments.
	$(DOCKER_RUN) export -c quote

.PHONY: run-export-all
run-export-all: ## Run export app to get all kbs.
	$(DOCKER_RUN) export

.PHONY: run-export-with-ns-cat
run-export-with-ns-cat: ## Run export app to get all kbs that match specific category and namespace.
	$(DOCKER_RUN) export -c quote -n Default

.PHONY: run-sync-with-args
run-sync-with-args: ## Run sync app with predefined arguments.
	$(DOCKER_RUN) sync --show-added-kbs --show-failed-kbs
