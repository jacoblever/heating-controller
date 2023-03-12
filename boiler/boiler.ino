#include <ESP8266WiFi.h>
#include <ESP8266HTTPClient.h>
#include <Stepper.h>
#include "Environment.h"

#define LED_PIN D3
#define STEPS 300

const char* ssid = WIFI_SSID;
const char* password = WIFI_PASSWORD;

WiFiClient wifiClient;

Stepper stepper(STEPS, D5, D7, D6, D8);

void setup() {
  pinMode(LED_PIN, OUTPUT);

  Serial.begin(9600);

  WiFi.begin(ssid, password);
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }
  Serial.print("Connected, IP address: ");
  Serial.println(WiFi.localIP());

  stepper.setSpeed(60);
}

String boilerState = "";

void loop() {
	if(WiFi.status() == WL_CONNECTED) {
    HTTPClient http;

    String serverPath = "http://192.168.86.100:8080/boiler-state/";

    http.begin(wifiClient, serverPath);
    
    int httpResponseCode = http.GET();

    if(httpResponseCode > 0) {
      Serial.print("HTTP Response code: ");
      Serial.println(httpResponseCode);
      String payload = http.getString();
      if(payload == "on") {
        digitalWrite(LED_PIN, HIGH);
        if(boilerState == "off") {
          Serial.println("Turning boiler off...");
          stepper.step(STEPS);
        }
      } else if(payload == "off") {
        digitalWrite(LED_PIN, LOW); 
        if(boilerState == "on") {
          Serial.println("Turning boiler on...");
          stepper.step(-STEPS);
        }     
      } else {
        Serial.print("Unknown boiler state: ");
        Serial.println(payload);
      }
      boilerState = payload;
    } else {
      Serial.print("Error code: ");
      Serial.println(httpResponseCode);
    }
    
    http.end();
  }
  
  delay(1000);
}
