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
	Spark.function("digitalread", sparkDigitalRead);
	Spark.function("digitalwrite", sparkDigitalWrite);

	Spark.function("analogread", sparkAnalogRead);
	Spark.function("analogwrite", sparkAnalogWrite);
}

int sparkDigitalRead(String pin)
{
	int pinNumber = pin.charAt(1) - '0';
	if (pinNumber < 0 || pinNumber > 7) { return -1; }

	if (pin.startsWith("D")) {
		pinMode(pinNumber, INPUT_PULLDOWN);
		return digitalRead(pinNumber);
	} else if (pin.startsWith("A")) {
		pinMode(pinNumber + 10, INPUT_PULLDOWN);
		return digitalRead(pinNumber + 10);
	}

	return -2;
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

	if (command.startsWith("D")) {
		pinMode(pinNumber, OUTPUT);
		digitalWrite(pinNumber, value);
		return 1;
	} else if (command.startsWith("A")) {
		pinMode(pinNumber + 10, OUTPUT);
		digitalWrite(pinNumber + 10, value);
		return 1;
	}

	return -3;
}

int sparkAnalogRead(String pin)
{
	int pinNumber = pin.charAt(1) - '0';
	if (pinNumber < 0 || pinNumber > 7) { return -1; }

	if (pin.startsWith("D")) {
		pinMode(pinNumber, INPUT);
		return analogRead(pinNumber);
	} else if (pin.startsWith("A")) {
		pinMode(pinNumber + 10, INPUT);
		return analogRead(pinNumber + 10);
	}

	return -2;
}

int sparkAnalogWrite(String command)
{
	int pinNumber = command.charAt(1) - '0';
	if (pinNumber < 0 || pinNumber > 7) { return -1; }

	String value = command.substring(3);

	if (command.startsWith("D")) {
		pinMode(pinNumber, OUTPUT);
		analogWrite(pinNumber, value.toInt());
		return 1;
	} else if (command.startsWith("A")) {
		pinMode(pinNumber + 10, OUTPUT);
		analogWrite(pinNumber + 10, value.toInt());
		return 1;
	}
	
	return -2;
}
