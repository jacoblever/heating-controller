.PHONY: connect
connect:
	ssh brain@192.168.86.100

.PHONY: deploy
deploy:
	cd brain && make deploy
	cd dashboard && make deploy
