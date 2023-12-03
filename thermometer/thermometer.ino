
#include <ESP8266WiFi.h>
#include <ESP8266HTTPClient.h>
#include <TsicSensor.h>
#include <SPI.h>
#include <Wire.h>
#include "Environment.h"

#define TEMP_PIN D5
#define DEFAULT_DELAY 1000

const char* ssid = WIFI_SSID;
const char* password = WIFI_PASSWORD;

WiFiClient wifiClient;

float temperature;
TsicSensor* sensor;

void setup() {

  Serial.begin(9600);

  WiFi.begin(ssid, password);
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }
  Serial.print("Connected, IP address: ");
  Serial.println(WiFi.localIP());

  // This creates/initializes a TSIC_506 sensor connected to GPIO 16.
  // -----------------------------------------------------------------------------------
  // The sensor is configured with "TsicExternalVcc", so it has a permanent external 
  // power source. The sensor values are then read in the background. 
  // (we can check for new values with the "newValueAvailable()" function...)
  sensor = TsicSensor::create(TEMP_PIN, TsicExternalVcc, TsicType::TSIC_306);
}


void loop() {
	if(sensor->newValueAvailable()) {      
    temperature = sensor->getTempCelsius();
  }

  Serial.println(temperature);
  
  if(temperature != 0) {
    unsigned long nextDelay = sendTemperature();
    delay(nextDelay);
    return;
  }

  delay(DEFAULT_DELAY);
}

unsigned long sendTemperature() {
  if(WiFi.status() != WL_CONNECTED) {
    return DEFAULT_DELAY;
  }

  HTTPClient http;

  String serverPath = "http://192.168.86.100:8080/temperature/?temperature=" + String(temperature) + "&id=2";
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

  return DEFAULT_DELAY;
}
