// Which pin on the Arduino is connected to the NeoPixels?
#define PIN        6 // On Trinket or Gemma, suggest changing this to 1

// How many NeoPixels are attached to the Arduino?
#define NUMPIXELS 3 // Popular NeoPixel ring size

#include <Adafruit_NeoPixel.h>
#ifdef __AVR__
 #include <avr/power.h> // Required for 16 MHz Adafruit Trinket
#endif

Adafruit_NeoPixel pixels(NUMPIXELS, PIN, NEO_GRB + NEO_KHZ800);

void setup() {
  // put your setup code here, to run once:
  Serial.begin(115200);
  pixels.begin();
}

uint8_t buff[3 * NUMPIXELS];
uint8_t buffPos = 0;

long lastRx = 0;

void loop() {
  while (Serial.available()) {
    buff[buffPos] = Serial.read();
    buffPos++;
    if (buffPos == 3 * NUMPIXELS) {
      pixels.clear();
      for (int i = 0; i < NUMPIXELS; i++) {
        pixels.setPixelColor(i, pixels.Color(buff[i*3], buff[i*3+1], buff[i*3+2]));
      }
      pixels.show();
      buffPos = 0;
    }
    lastRx = millis();
  }
  if (millis() - lastRx > 500) {
    buffPos = 0;
  }
}
