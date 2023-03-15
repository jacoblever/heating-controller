# dashboard

Displays current state of the system. This is a simple static webpage served using python.

## Making the server auto run

Following https://forums.raspberrypi.com//viewtopic.php?p=921354

Copy the service config onto the pi
```
scp dashboard.service brain@192.168.86.100:/home/brain
```

SSH into brain, move it and give it the correct permitions
```
sudo cp dashboard.service /etc/systemd/system/
sudo chmod u+rwx /etc/systemd/system/dashboard.service
```

Enable it and start it (or you could just reboot)
```
sudo systemctl enable dashboard
sudo systemctl start dashboard
```
