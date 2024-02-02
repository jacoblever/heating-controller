# heating-controller

## Arduino project dependencies

- In the Arduino IDE menu click on File → Preferences
- Insert the following URL in the field “Additional Boards Manager URLs:”: http://arduino.esp8266.com/stable/package_esp8266com_index.json (was https://raw.githubusercontent.com/espressif/arduino-esp32/gh-pages/package_esp32_index.json)
- Click in the Arduino IDE on Tools → Board → Board Manager.
- Now search for NodeMCU or ESP8266 and you will find the esp8266 by ESP8266 Community. Install the latest version of the board.
- Click in the Arduino IDE on Tools → Board → esp8266. Select the ESP8266 board with the name "NodeMCU 1.0 (ESP-12E Module)"

These steps come from following the instructions in Section "How to flash your Code on the ESP8266 WeMos D1 Mini" of this tutorial: 
https://diyi0t.com/esp8266-wemos-d1-mini-tutorial/

Apple Silicon Macs need installation of Rosetta for the compilation to work. This just involves running `softwareupdate --install-rosetta` in Terminal and accepting the license agreement. More details can be found: here https://support.arduino.cc/hc/en-us/articles/7765785712156-Error-bad-CPU-type-in-executable-on-macOS


Libraries used by Arduino project:

* `Arduino_JSON` by `Arduino`
* https://github.com/tripplefox/TsicSensor
* https://github.com/olikraus/u8g2
* https://github.com/plageoj/urlencode
