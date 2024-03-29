# Setup Rasberry Pi

Config Assumptions
```
IP address of raspberry pi: 192.168.86.100
Username of raspberry pi user account: brain
```

## Passwordless ssh

Following https://levelup.gitconnected.com/how-to-connect-without-password-using-ssh-passwordless-9b8963c828e8 to setup passwordless login

Assuming you already have a local pub/privaite key pair generated on your local machine, just run

```
ssh-copy-id brain@192.168.86.100
```

## Making the server auto run

Following https://forums.raspberrypi.com//viewtopic.php?p=921354

Copy the service config onto the pi
```
scp server.service brain@192.168.86.100:/home/brain
```

Copy the env file onto the pi (make sure you have created a `.env` file, based on `.env.template`, with the correct values)

```
scp .env brain@192.168.86.100:/home/brain/.env
```

SSH into brain, move it and give it the correct permitions
```
sudo cp server.service /etc/systemd/system/
sudo chmod u+rwx /etc/systemd/system/server.service
```

Enable it and start it (or you could just reboot)
```
sudo systemctl enable server
sudo systemctl start server
```

# Compiling the go server for Raspberry Pi

See https://dev.to/coreyvan/from-zero-to-http-servers-with-go-and-raspberry-pi-3oi1

```
GOOS=linux GOARCH=arm GOARM=7 go build
```