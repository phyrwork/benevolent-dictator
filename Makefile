.PHONY: api

api:
	make -C pkg/api
	make -C cmd/api
