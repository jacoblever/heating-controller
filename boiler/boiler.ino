#include <ESP8266WiFi.h>
#include <ESP8266HTTPClient.h>
#include <Stepper.h>
#include <Arduino_JSON.h>
#include <UrlEncode.h>
#include "Environment.h"

#define ERROR_LED_PIN D3
#define BOILER_STATUS_LED_PIN D2
#define STEPS_IN_ONE_REVOLUTION 4096

#define LOG_BUFFER_SIZE 10

WiFiClient wifiClient;

Stepper stepper(STEPS_IN_ONE_REVOLUTION, D5, D7, D6, D8);

int pollDelayMs = 1000;
String boilerState = "";
bool inError = false;

String logger[LOG_BUFFER_SIZE];
int nextLoggerIndex = 0;

void setup() {
  Serial.begin(9600);
  pinMode(ERROR_LED_PIN, OUTPUT);
  pinMode(BOILER_STATUS_LED_PIN, OUTPUT);
  digitalWrite(ERROR_LED_PIN, HIGH);

  connectWiFi(WIFI_SSID, WIFI_PASSWORD);

  logToServer("Arduino setup complete");
}

void loop() {
	if(WiFi.status() == WL_CONNECTED) {
    JSONVar response = makeRequest("http://192.168.86.100:8080/boiler-state/");
    if(response.hasOwnProperty("PollDelayMs")) {
      pollDelayMs = response["PollDelayMs"];
    }

    if(int(response["HTTPStatus"]) != 200) {
      setErrorMode(true);
      delay(pollDelayMs);
      return;
    }

    String newBoilerState = response["BoilerState"];
    int motorSpeedRpm = response["MotorSpeedRPM"];
    int stepsToTurn = response["StepsToTurn"];
    String command = response["Command"];

    if(newBoilerState == "on") {
      digitalWrite(BOILER_STATUS_LED_PIN, HIGH);

      if(boilerState == "off") {
        logToServer("Turning boiler on...");
        turnMotor(stepper, stepsToTurn, motorSpeedRpm);
      }
    } else if(newBoilerState == "off") {
      digitalWrite(BOILER_STATUS_LED_PIN, LOW);

      if(boilerState == "on") {
        logToServer("Turning boiler off...");
        turnMotor(stepper, -stepsToTurn, motorSpeedRpm);
      }     
    } else {
      logToServer("Unknown boiler state: " + newBoilerState);
    }
    boilerState = newBoilerState;

    if (newBoilerState == "on" || newBoilerState == "off") {
      setErrorMode(false);
    }

    if(command != ""){
      int steps = command.toInt();
      logToServer("Running command: " + command);
      turnMotor(stepper, steps, motorSpeedRpm);
	  }
  }
  
  delay(pollDelayMs);
}

void setErrorMode(bool on) {
  if(on) {
    if(!inError) {
      logToServer("Connection to server lost");
    }
    digitalWrite(ERROR_LED_PIN, HIGH);
  } else {
    if(inError) {
      logToServer("Connection to server restored");
    }
    digitalWrite(ERROR_LED_PIN, LOW);
  }
  inError = on;
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
  
  Serial.println();
  logToServer("Connected to WiFi with IP address: " + WiFi.localIP().toString());
}

JSONVar makeRequest(String url) {
  HTTPClient http;

  http.begin(wifiClient, url);
  http.addHeader("Content-Type", "application/x-www-form-urlencoded");

  Serial.print("Making HTTP request to: ");
  Serial.println(url);

  int httpResponseCode = http.POST(getServerLogsPostData());

  String body = "{}";
  if (httpResponseCode > 0) {
    body = http.getString();
    String outcomeLog = "Status: " + String(httpResponseCode) + ", Body: " + body;

    if(httpResponseCode == 200) {
      clearServerLogs();
      Serial.println(outcomeLog);
    } else {
      logToServer(outcomeLog);
    }
  } else {
    logToServer("Status: " + String(httpResponseCode));
  }

  http.end();

  JSONVar response = JSON.parse(body);
  if (JSON.typeof(response) == "undefined") {
    logToServer("Parsing response failed!");
  }
  response["HTTPStatus"] = httpResponseCode;
  return response;
}

void logToServer(String message) {
  Serial.println(message);

  if (nextLoggerIndex == LOG_BUFFER_SIZE - 1) {
    String reachedMaxMsg = "Server log buffer full. This means some genuine logs were lost. Please increase LOG_BUFFER_SIZE to capture full logs.";
    logger[nextLoggerIndex] = reachedMaxMsg;
    Serial.println(reachedMaxMsg);
    return;
  }

  logger[nextLoggerIndex] = message;
  nextLoggerIndex = nextLoggerIndex + 1;
}

String getServerLogsPostData() {
  String logs = "";

  for (int i = 0; i < LOG_BUFFER_SIZE; i++) {
    if (logger[i] != "") {
      if(logs != "") {
        logs += "&";
      }
      logs += "Log=" + urlEncode(logger[i]);
    }
  }
  return logs;
}

void clearServerLogs() {
  for (int i = 0; i < LOG_BUFFER_SIZE; i++) {
    if (logger[i] != "") {
      logger[i] = "";
    }
  }
  nextLoggerIndex = 0;
}
