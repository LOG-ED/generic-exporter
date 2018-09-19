.PHONY: build start stop help

help:
	@echo "Please use \`make <target>' where <target> is one of"
	@echo "  build			when golang files are modified"
	@echo "  start     		to start the stack"
	@echo "  stop			to destroy the stack"

build:
	docker-compose build --no-cache

start:
	docker-compose up -d --remove-orphans

stop: 
	docker-compose down --remove-orphans