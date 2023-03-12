#include <ESP8266WiFi.h>
#include <ESP8266HTTPClient.h>
#include <TsicSensor.h>
#include <U8g2lib.h>
#include <SPI.h>
#include <Wire.h>
#include "Environment.h"

#define LED_PIN D3
#define TEMP_PIN D5

const char* ssid = WIFI_SSID;
const char* password = WIFI_PASSWORD;

WiFiClient wifiClient;

float temperature;
TsicSensor* sensor;

U8G2_SH1107_SEEED_128X128_1_SW_I2C u8g2(U8G2_R0, /* clock=*/ D1, /* data=*/ D2, /* reset=*/ U8X8_PIN_NONE);

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

  // This creates/initializes a TSIC_506 sensor connected to GPIO 16.
  // -----------------------------------------------------------------------------------
  // The sensor is configured with "TsicExternalVcc", so it has a permanent external 
  // power source. The sensor values are then read in the background. 
  // (we can check for new values with the "newValueAvailable()" function...)
  sensor = TsicSensor::create(TEMP_PIN, TsicExternalVcc, TsicType::TSIC_306);

  u8g2.begin();
}


void loop() {
	if(sensor->newValueAvailable()) {      
   temperature = sensor->getTempCelsius();
  }

  if(temperature < 24) {
    digitalWrite(LED_PIN, HIGH);
  } else {
    digitalWrite(LED_PIN, LOW);      
  }

  u8g2.firstPage();
  do {
    u8g2.setFont(u8g2_font_ncenB10_tr);
    u8g2.drawStr(0, 24, String(temperature).c_str());
  } while ( u8g2.nextPage() );

  Serial.println(temperature);
  
  if(WiFi.status() == WL_CONNECTED) {
    HTTPClient http;

    String serverPath = "http://192.168.86.100:8080/update-temperature/?temperature=" + String(temperature);
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