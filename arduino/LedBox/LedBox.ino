#include <Adafruit_NeoPixel.h>

#include <Adafruit_GFX.h>
#include <Adafruit_NeoMatrix.h>
#include <Adafruit_NeoPixel.h>
#ifndef PSTR
 #define PSTR
#endif

#include <ESP8266WiFi.h>
#include <ArduinoJson.h>

#define STRIP_PIN       D2
#define MATRIX_PIN      D1
#define STRIP_N         60
#define MATRIX_N        16

Adafruit_NeoPixel strip(STRIP_N, STRIP_PIN, NEO_GRB + NEO_KHZ800);
Adafruit_NeoMatrix matrix = Adafruit_NeoMatrix(MATRIX_N, MATRIX_N, MATRIX_PIN,
  NEO_MATRIX_TOP     + NEO_MATRIX_LEFT +
  NEO_MATRIX_COLUMNS + NEO_MATRIX_ZIGZAG,
  NEO_GRB            + NEO_KHZ800);

const char* ssid     = "w";
const char* password = "fsociety";
const char* host = "c30";
const uint16_t port = 8080;

StaticJsonDocument<200> doc;

const uint16_t colors[] = {
  matrix.Color(255, 0, 0), matrix.Color(0, 255, 0), matrix.Color(0, 0, 255) };

void indicate(int r, int g, int b) {
  matrix.drawPixel(0,0, matrix.Color(r,g,b));
  matrix.show();
}

struct Indicators {
  bool mkr_1010;
};

struct Indicators indicators;
WiFiClient client;

const uint16_t indicatorColors[2] = {
  matrix.Color(255,0,0),
  matrix.Color(0,255,0),
};

void updateIndicators() {
  matrix.clear();
  matrix.drawPixel(0,0,indicatorColors[WiFi.status() == WL_CONNECTED ? 1 : 0]);
  matrix.drawPixel(1,0,indicatorColors[client.connected() ? 1 : 0]);
  matrix.drawPixel(3,0,indicatorColors[1]);

  matrix.show();
}

void setup() {
  Serial.begin(9600);

  strip.begin();
  matrix.begin();
  matrix.setTextWrap(false);
  matrix.setBrightness(40);
  matrix.setTextColor(colors[0]);

  updateIndicators();

  Serial.print("connecting to wifi");
  WiFi.mode(WIFI_STA);
  WiFi.begin(ssid, password);
  while (WiFi.status() != WL_CONNECTED) {
    Serial.print(".");
    delay(500);
  }
  Serial.println("\nwifi connected");
  updateIndicators();
}

int x    = matrix.width();
int pass = 0;

void loop() {
  updateIndicators();

  if (!client.connect(host, port)) {
    Serial.println("failed to connect to wifi");
    delay(5000);
    return;
  }

  updateIndicators();

  if (client.connected()) {
    client.println("register dom_leds doorbell,:active");
  }
  while(true) {
    String lineStr = client.readStringUntil('\n');
    if (lineStr.length() > 0) {
      char line[lineStr.length()];
      lineStr.toCharArray(line, lineStr.length() + 1);

      DeserializationError error = deserializeJson(doc, line);
      if (error) {
        indicate(255,0,0);
        delay(1500);
        continue;
      }
      const char *event = doc["event"];
      Serial.print("Event: ");
      Serial.println(event);
      delay(500);
      if (strcmp(event, "doorbell")==0) {
        Serial.println("ff");
        delay(500);
      } else if (strcmp(event, ":active") ==0) {
        char *arr[] = doc["payload"];
        char *s = "mkr_1010";
        int i = 0;
        while (arr[i]) {
          if(strcmp(arr[i], s) == 0) {
            indicators.mkr_1010 = true;
            break;
          }
          i++;
        }
      }
    }
    delay(50);
  }


  delay(50000);
}
