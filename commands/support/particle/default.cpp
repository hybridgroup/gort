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

  pinMode(pinNumber, OUTPUT);
  digitalWrite(pinNumber, value);

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
  int pinNum = command.charAt(1) - '0';
  char cmdType = command.charAt(0);

  if (pinNum < 0 || pinNum > 9) { return -1; }

  String value = command.substring(3);

  if (cmdType == 'D') {
    pinMode(pinNum, OUTPUT);
    analogWrite(pinNum, value.toInt());
    return 1;
  } else if (cmdType == 'A') {
    pinMode(pinNum + 10, OUTPUT);
    analogWrite(pinNum + 10, value.toInt());
    return 2;
  } else if (cmdType == 'S') {
    if (pinNum < 8) {
      pinNum += 10;
    } else if (pinNum > 7) {
      pinNum -= 8;
    }
    Servo myServo;
    myServo.attach(pinNum);
    myServo.write(value.toInt());
    return 3;
  }

  return -2;
}
