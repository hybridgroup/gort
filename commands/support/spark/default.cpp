/**
 * Spark Core - THG Default Firmware
 *
 * Copyright (c) 2014 The Hybrid Group
 * Licensed under the Apache 2.0 license.
 */

// Includes
#include "application.h"

// Function prototypes
int sparkDigitalRead(String pin);
int sparkDigitalWrite(String command);
int sparkAnalogRead(String pin);
int sparkAnalogWrite(String command);

void setup()
{
  Spark.function("digitalwrite", sparkDigitalWrite);
  Spark.function("digitalread", sparkDigitalRead);

  Spark.function("analogwrite", sparkAnalogWrite);
    Spark.function("analogread", sparkAnalogRead);
}

int sparkDigitalRead(String pin)
{
  int pinNumber = pin.charAt(1) - '0';

  if (pinNumber < 0 || pinNumber > 7) { return -1; }
  if (pin.startsWith("A")) { pinNumber += 10; }

  pinMode(pinNumber, INPUT_PULLDOWN);

  return digitalRead(pinNumber);
}

int sparkDigitalWrite(String command)
{
  bool value = 0;
  int pinNumber = command.charAt(1) - '0';

  if (pinNumber < 0 || pinNumber > 7) { return -1; }

  if (command.substring(3, 7) == "HIGH") {
    value = 1;
  } else if (command.substring(3, 6) == "LOW") {
    value = 0;
  } else {
    return -2;
  }

  if (command.startsWith("A")) { pinNumber += 10; }

  pinMode(pinNumber + 10, OUTPUT);
  digitalWrite(pinNumber + 10, value);

  return 1;
}

int sparkAnalogRead(String pin)
{
  int pinNumber = pin.charAt(1) - '0';

  if (pinNumber < 0 || pinNumber > 7) { return -1; }
  if (pin.startsWith("A")) { pinNumber += 10; }

  pinMode(pinNumber, INPUT);

  return analogRead(pinNumber);
}

int sparkAnalogWrite(String command)
{
  Servo myServo;
  bool servo = false;
  int pinNumber = command.charAt(1) - '0';

  if (pinNumber < 0 || pinNumber > 7) { return -1; }

  if (toupper(command.charAt(3)) == 'S') { servo = true; }

  String value = command.substring(5);

  if (command.startsWith("A")) { pinNumber += 10; }

  if (servo) {
    myServo.attach(pinNumber);
    myServo.write(value.toInt());
    return 1;
  } else {
    pinMode(pinNumber, OUTPUT);
    analogWrite(pinNumber, value.toInt());
    return 1;
  }

  return -2;
}
