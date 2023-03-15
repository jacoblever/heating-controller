#include <ESP8266WiFi.h>
#include <ESP8266HTTPClient.h>
#include <Arduino_JSON.h>
#include "Environment.h"

WiFiClient wifiClient;

int pollDelayMs = 1000;

void setup() {
  Serial.begin(9600);
  connectWiFi(WIFI_SSID, WIFI_PASSWORD);
}

void loop() { 
  if(WiFi.status() == WL_CONNECTED) {
    JSONVar response = makeRequest("http://192.168.86.100:8080/smart-switch-alive/");
    if(response.hasOwnProperty("PollDelayMs")) {
      pollDelayMs = response["PollDelayMs"];
    }
  }
  delay(pollDelayMs);
}

void connectWiFi(const char* ssid, const char* password) {
  Serial.print("Connecting to WiFi network");
  WiFi.begin(ssid, password);
  while (WiFi.status() != WL_CONNECTED) {
    delay(200);
    Serial.print(".");
  }
  
  Serial.print("\nConnected, IP address: ");
  Serial.println(WiFi.localIP());
}

JSONVar makeRequest(String url) {
  HTTPClient http;

  http.begin(wifiClient, url);

  Serial.print("Making HTTP request to: ");
  Serial.println(url);

  int httpResponseCode = http.GET();

  String body = "{}"; 
  if (httpResponseCode > 0) {
    Serial.print("Status: ");
    Serial.print(httpResponseCode);
    body = http.getString();
    Serial.print(" - ");
    Serial.println(body);
  }
  else {
    Serial.print("Error code: ");
    Serial.print(httpResponseCode);
  }

  http.end();

  JSONVar response = JSON.parse(body);
  if (JSON.typeof(response) == "undefined") {
    Serial.println("Parsing input failed!");
  }
  return response;
}
