.PHONY: connect
connect:
	ssh brain@192.168.86.100

.PHONY: deploy
deploy:
	npm run build
	ssh brain@192.168.86.100 'rm -rf /home/brain/dashboard/*'
	scp -r ./build/* brain@192.168.86.100:/home/brain/dashboard/
