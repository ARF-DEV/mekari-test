.PHONY: docs

docs:
	swag fmt
	swag init -g api/api.go