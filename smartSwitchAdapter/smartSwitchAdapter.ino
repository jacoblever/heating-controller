#include <ESP8266WiFi.h>
#include <ESP8266HTTPClient.h>
#include "Environment.h"

const char* ssid = WIFI_SSID;
const char* password = WIFI_PASSWORD;

WiFiClient wifiClient;

void setup() {
  Serial.begin(9600);

  WiFi.begin(ssid, password);
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }
  
  Serial.print("Connected, IP address: ");
  Serial.println(WiFi.localIP());
}

void loop() { 
  if(WiFi.status() == WL_CONNECTED) {
    HTTPClient http;

    String serverPath = "http://192.168.86.100:8080/smart-switch-alive/";

    http.begin(wifiClient, serverPath);
    int httpResponseCode = http.GET();

    if(httpResponseCode > 0) {
      Serial.print("HTTP Response code: ");
      Serial.println(httpResponseCode);
      String payload = http.getString();
      Serial.println(payload);
    } else {
      Serial.print("Error code: ");
      Serial.println(httpResponseCode);
    }

    http.end();
  }

  delay(1000);
}
