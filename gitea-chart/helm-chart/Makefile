.PHONY: prepare-environment
prepare-environment:
	npm install

.PHONY: readme
readme: prepare-environment
	npm run readme:parameters
	npm run readme:lint

.PHONY: unittests
unittests:
	helm unittest --helm3 --strict -f 'unittests/**/*.yaml' ./
