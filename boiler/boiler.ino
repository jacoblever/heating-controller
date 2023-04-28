#include <ESP8266WiFi.h>
#include <ESP8266HTTPClient.h>
#include <Stepper.h>
#include <Arduino_JSON.h>
#include "Environment.h"

#define STANDBY_LED_PIN D3
#define BOILER_STATUS_LED_PIN D2
#define STEPS_IN_ONE_REVOLUTION 4096

WiFiClient wifiClient;

Stepper stepper(STEPS_IN_ONE_REVOLUTION, D5, D7, D6, D8);

int pollDelayMs = 1000;
String boilerState = "";

void setup() {
  Serial.begin(9600);
  pinMode(STANDBY_LED_PIN, OUTPUT);
  pinMode(BOILER_STATUS_LED_PIN, OUTPUT);
  digitalWrite(STANDBY_LED_PIN, LOW);

  connectWiFi(WIFI_SSID, WIFI_PASSWORD);
}

void loop() {
	if(WiFi.status() == WL_CONNECTED) {
    JSONVar response = makeRequest("http://192.168.86.100:8080/boiler-state/");
    if(response.hasOwnProperty("PollDelayMs")) {
      pollDelayMs = response["PollDelayMs"];
    }

    String newBoilerState = response["BoilerState"];
    int motorSpeedRpm = response["MotorSpeedRPM"];
    int stepsToTurn = response["StepsToTurn"];
    String command = response["Command"];

    if(newBoilerState == "on") {
      digitalWrite(BOILER_STATUS_LED_PIN, HIGH);

      if(boilerState == "off") {
        Serial.println("Turning boiler on...");
        turnMotor(stepper, stepsToTurn, motorSpeedRpm);
      }
    } else if(newBoilerState == "off") {
      digitalWrite(BOILER_STATUS_LED_PIN, LOW);

      if(boilerState == "on") {
        Serial.println("Turning boiler off...");
        turnMotor(stepper, -stepsToTurn, motorSpeedRpm);
      }     
    } else {
      Serial.print("Unknown boiler state: ");
      Serial.println(newBoilerState);
    }
    boilerState = newBoilerState;

    if (boilerState == "on" || boilerState == "off") {
      digitalWrite(STANDBY_LED_PIN, HIGH);
    }

    if(command != ""){
      int steps = command.toInt();
      Serial.print("Running command: ");
      Serial.println(command);
      turnMotor(stepper, steps, motorSpeedRpm);
	  }
  }
  
  delay(pollDelayMs);
}

void turnMotor(Stepper stepper, int steps, int speedRpm) {
  Serial.print("Turning stepper motor ");
  Serial.print(steps);
  Serial.print(" steps at ");
  Serial.print(speedRpm);
  Serial.print(" RPM...");

  stepper.setSpeed(speedRpm);
  stepper.step(steps);

  Serial.println(" Done!");
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
