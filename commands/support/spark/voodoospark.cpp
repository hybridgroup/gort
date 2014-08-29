/**
  ******************************************************************************
  * @file    voodoospark.cpp
  * @author  Chris Williams
  * @version V2.3.1
  * @date    1-June-2014
  * @brief   Exposes the firmware level API through a TCP Connection initiated
  *          to the spark device
  ******************************************************************************
  Copyright (c) 2014 Chris Williams  All rights reserved.

  Permission is hereby granted, free of charge, to any person
  obtaining a copy of this software and associated documentation
  files (the "Software"), to deal in the Software without
  restriction, including without limitation the rights to use,
  copy, modify, merge, publish, distribute, sublicense, and/or sell
  copies of the Software, and to permit persons to whom the
  Software is furnished to do so, subject to the following
  conditions:

  The above copyright notice and this permission notice shall be
  included in all copies or substantial portions of the Software.

  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
  EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
  OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
  NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
  HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
  WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
  FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
  OTHER DEALINGS IN THE SOFTWARE.
  ******************************************************************************
  */
#include "application.h"

#define DEBUG 0

// Port = 0xbeef
#define PORT 48879

// allow use of itoa() in this scope
extern char* itoa(int a, char* buffer, unsigned char radix);

// table of action codes
// to do: make this an enum?
#define msg_pinMode                    (0x00)
#define msg_digitalWrite               (0x01)
#define msg_analogWrite                (0x02)
#define msg_digitalRead                (0x03)
#define msg_analogRead                 (0x04)
#define msg_setAlwaysSendBit           (0x05)
#define msg_setSampleInterval          (0x06)
/* NOTE GAP */
// #define msg_serialBegin                (0x10)
// #define msg_serialEnd                  (0x11)
// #define msg_serialPeek                 (0x12)
// #define msg_serialAvailable            (0x13)
// #define msg_serialWrite                (0x14)
// #define msg_serialRead                 (0x15)
// #define msg_serialFlush                (0x16)
/* NOTE GAP */
// #define msg_spiBegin                   (0x20)
// #define msg_spiEnd                     (0x21)
// #define msg_spiSetBitOrder             (0x22)
// #define msg_spiSetClockDivider         (0x23)
// #define msg_spiSetDataMode             (0x24)
// #define msg_spiTransfer                (0x25)
// /* NOTE GAP */
// #define msg_wireBegin                  (0x30)
// #define msg_wireRequestFrom            (0x31)
// #define msg_wireBeginTransmission      (0x32)
// #define msg_wireEndTransmission        (0x33)
// #define msg_wireWrite                  (0x34)
// #define msg_wireAvailable              (0x35)
// #define msg_wireRead                   (0x36)
/* NOTE GAP */
#define msg_servoWrite                 (0x41)

#define msg_count                      (0x46)

uint8_t bytesPerAction[] = {
  // digital/analog I/O
  2,    // msg_pinMode
  2,    // msg_digitalWrite
  2,    // msg_analogWrite
  1,    // msg_digitalRead
  1,    // msg_analogRead
  2,    // msg_setAlwaysSendBit
  1,    // msg_setSampleInterval
  // gap from 0x07-0x0f
  0,    // msg_0x07
  0,    // msg_0x08
  0,    // msg_0x09
  0,    // msg_0x0a
  0,    // msg_0x0b
  0,    // msg_0x0c
  0,    // msg_0x0d
  0,    // msg_0x0e
  0,    // msg_0x0f
  // serial I/O
  2,    // msg_serialBegin
  1,    // msg_serialEnd
  1,    // msg_serialPeek
  1,    // msg_serialAvailable
  2,    // msg_serialWrite  -- variable length message!
  1,    // msg_serialRead
  1,    // msg_serialFlush
  // gap from 0x17-0x1f
  0,    // msg_0x17
  0,    // msg_0x18
  0,    // msg_0x19
  0,    // msg_0x1a
  0,    // msg_0x1b
  0,    // msg_0x1c
  0,    // msg_0x1d
  0,    // msg_0x1e
  0,    // msg_0x1f
  // SPI I/O
  0,    // msg_spiBegin
  0,    // msg_spiEnd
  1,    // msg_spiSetBitOrder
  1,    // msg_spiSetClockDivider
  1,    // msg_spiSetDataMode
  1,    // msg_spiTransfer
  // gap from 0x26-0x2f
  0,    // msg_0x26
  0,    // msg_0x27
  0,    // msg_0x28
  0,    // msg_0x29
  0,    // msg_0x2a
  0,    // msg_0x2b
  0,    // msg_0x2c
  0,    // msg_0x2d
  0,    // msg_0x2e
  0,    // msg_0x2f
  // wire I/O
  1,    // msg_wireBegin
  3,    // msg_wireRequestFrom
  1,    // msg_wireBeginTransmission
  1,    // msg_wireEndTransmission
  1,    // msg_wireWrite  -- variable length message!
  0,    // msg_wireAvailable
  0,    // msg_wireRead
  // gap from 0x37-0x3f
  0,    // msg_0x37
  0,    // msg_0x38
  0,    // msg_0x39
  0,    // msg_0x3a
  0,    // msg_0x3b
  0,    // msg_0x3c
  0,    // msg_0x3d
  0,    // msg_0x3e
  0,    // msg_0x3f
  0,    // msg_0x40
  // servo
  2,    // msg_servoWrite
  1,    // msg_servoDetach
};


TCPServer server = TCPServer(PORT);
TCPClient client;

bool hasAction = false;
bool isConnected = false;

byte reporting[20];
byte buffer[16];
byte cached[4];

int reporters = 0;
int bytesRead = 0;
int bytesExpecting = 0;
int action, available;

unsigned long lastms;
unsigned long nowms;
unsigned long sampleInterval = 100;
unsigned long SerialSpeed[] = {
  600, 1200, 2400, 4800, 9600, 14400, 19200, 28800, 38400, 57600, 115200
};





/*
  PWM/Servo support is CONFIRMED available on:

  D0, D1, A0, A1, A4, A5, A6, A7

  Allocate 8 servo objects:
 */
Servo servos[8];
/*
  The Spark board can only support PWM/Servo on specific pins, so
  based on the pin number, determine the servo index for the allocated
  servo object.
 */
int ToServoIndex(int pin) {
  // D0, D1
  if (pin == 0 || pin == 1) return pin;
  // A0, A1
  if (pin == 10 || pin == 11) return pin - 8;
  // A4, A5, A6, A7
  if (pin >= 14) return pin - 10;
}

void send(int action, int pin, int value) {
  // See https://github.com/voodootikigod/voodoospark/issues/20
  // to understand why the send function splits values
  // into two 7-bit bytes before sending.
  //
  int lsb = value & 0x7f;
  int msb = value >> 0x07 & 0x7f;

  server.write(action);
  server.write(pin);

  // Send the LSB
  server.write(lsb);
  // Send the MSB
  server.write(msb);

  // #ifdef DEBUG
  // Serial.print("SENT: ");
  // Serial.print(value);
  // Serial.print(" -> [ ");
  // Serial.print(lsb);
  // Serial.print(", ");
  // Serial.print(msb);
  // Serial.println(" ]");
  // #endif
}

void report() {
  if (isConnected) {
    for (int i = 0; i < 20; i++) {
      if (reporting[i]) {
        #ifdef DEBUG
        Serial.print("Reporting: ");
        Serial.print(i, DEC);
        Serial.println(reporting[i], DEC);
        #endif

        int dr = (reporting[i] & 1);
        int ar = (reporting[i] & 2);

        if (i < 10 && dr) {
          send(0x03, i, digitalRead(i));
        } else {
          if (dr) {
            send(0x03, i, digitalRead(i));
          } else {
            if (ar) {
              int adc = analogRead(i);
              #ifdef DEBUG
              Serial.print("Analog Report (pin, adc): ");
              Serial.print(i, DEC);
              Serial.println(adc, DEC);
              #endif
              send(0x04, i, adc);
            }
          }
        }
      }
    }
  }
}

void reset() {
  #ifdef DEBUG
  Serial.println("RESETTING");
  #endif

  bytesExpecting = 0;
  bytesRead = 0;
  hasAction = false;

  for (int i = 0; i < 20; i++) {
    // Clear the pin reporting list
    reporting[i] = 0;

    // Clear the incoming buffer
    if (i < 16) {
      buffer[i] = 0;
    }

    // Clear the action data cache
    if (i < 4) {
      buffer[i] = 0;
    }

    // Detach any attached servos
    if (i < 8) {
      if (servos[i].attached()) {
        servos[i].detach();
      }
    }
  }
}




void setup() {

  server.begin();
  netapp_ipconfig(&ip_config);

  #ifdef DEBUG
  Serial.begin(115200);
  #endif

  IPAddress myIp = WiFi.localIP();
  static char myIpString[24] = "";
  char octet[5];
  itoa(myIp[0],octet,10); strcat(myIpString,octet); strcat(myIpString,".");
  itoa(myIp[1],octet,10); strcat(myIpString,octet); strcat(myIpString,".");
  itoa(myIp[2],octet,10); strcat(myIpString,octet); strcat(myIpString,".");
  itoa(myIp[3],octet,10); strcat(myIpString,octet); strcat(myIpString,":");
  itoa(PORT,octet,10); strcat(myIpString,octet);
  Spark.variable("endpoint", myIpString, STRING);

}


void processInput() {
  int pin, mode, val, type, speed, address, stop, len, k, i;
  int byteCount = bytesRead;

  #ifdef DEBUG
  Serial.println("----------processInput----------");
  Serial.print("Bytes Available: ");
  Serial.println(available, DEC);

  for (i = 0; i < available; i++) {
    Serial.print(i, DEC);
    Serial.print(": ");
    Serial.println(buffer[i], DEC);
  }
  #endif

  // Only check if buffer[0] is possibly an action
  // when there is no known action in memory.
  if (hasAction == false) {
    if (buffer[0] < msg_count) {
      hasAction = true;
      bytesExpecting = bytesPerAction[action] + 1;
    }
  }

  #ifdef DEBUG
  Serial.print("Bytes Expecting: ");
  Serial.println(bytesExpecting, DEC);
  Serial.print("Bytes Read: ");
  Serial.println(bytesRead, DEC);
  #endif

  // When the first byte of buffer is an action and
  // enough bytes are read, begin processing the action.
  if (hasAction && bytesRead >= bytesExpecting) {

    action = buffer[0];

    #ifdef DEBUG
    Serial.print("Action received: ");
    Serial.println(action, DEC);
    #endif


    // Copy the expected bytes into the cache and shift
    // the unused bytes to the beginning of the buffer
    for (k = 0; k < byteCount; k++) {
      // Cache the bytes that we're expecting for
      // this action.
      if (k < bytesExpecting) {
        cached[k] = buffer[k];

        // Reduce the bytesRead by the number of bytes "taken"
        bytesRead--;

        #ifdef DEBUG
        Serial.print("Cached: ");
        Serial.print(k, DEC);
        Serial.println(cached[k], DEC);
        #endif
      }

      // Shift the unused buffer to the front
      buffer[k] = buffer[k + bytesExpecting];
    }

    // Proceed with action processing
    switch (action) {
      case msg_pinMode:  // pinMode
        pin = cached[1];
        mode = cached[2];
        #ifdef DEBUG
        Serial.print("PIN received: ");
        Serial.println(pin);
        Serial.print("MODE received: ");
        Serial.println(mode, HEX);
        #endif

        if (servos[ToServoIndex(pin)].attached()) {
          servos[ToServoIndex(pin)].detach();
        }

        if (mode == 0x00) {
          pinMode(pin, INPUT);
        } else if (mode == 0x02) {
          pinMode(pin, INPUT_PULLUP);
        } else if (mode == 0x03) {
          pinMode(pin, INPUT_PULLDOWN);
        } else if (mode == 0x01) {
          pinMode(pin, OUTPUT);
        } else if (mode == 0x04) {
          pinMode(pin, OUTPUT);
          servos[ToServoIndex(pin)].attach(pin);
        }
        break;

      case msg_digitalWrite:  // digitalWrite
        pin = cached[1];
        val = cached[2];
        #ifdef DEBUG
        Serial.print("PIN received: ");
        Serial.println(pin, DEC);
        Serial.print("VALUE received: ");
        Serial.println(val, HEX);
        #endif
        digitalWrite(pin, val);
        break;

      case msg_analogWrite:  // analogWrite
        pin = cached[1];
        val = cached[2];
        #ifdef DEBUG
        Serial.print("PIN received: ");
        Serial.println(pin, DEC);
        Serial.print("VALUE received: ");
        Serial.println(val, HEX);
        #endif
        analogWrite(pin, val);
        break;

      case msg_digitalRead:  // digitalRead
        pin = cached[1];
        val = digitalRead(pin);
        #ifdef DEBUG
        Serial.print("PIN received: ");
        Serial.println(pin, DEC);
        Serial.print("VALUE sent: ");
        Serial.println(val, HEX);
        #endif
        send(0x03, pin, val);
        break;

      case msg_analogRead:  // analogRead
        pin = cached[1];
        val = analogRead(pin);
        #ifdef DEBUG
        Serial.print("PIN received: ");
        Serial.println(pin, DEC);
        Serial.print("VALUE sent: ");
        Serial.println(val, HEX);
        #endif
        send(0x04, pin, val);
        break;

      case msg_setAlwaysSendBit: // set always send bit
        reporters++;
        pin = cached[1];
        val = cached[2];
        #ifdef DEBUG
        Serial.print("READ: ");
        Serial.print(pin, DEC);
        Serial.println(val, HEX);
        #endif
        reporting[pin] = val;
        break;

      case msg_setSampleInterval: // set the sampling interval in ms
        val = cached[1];
        sampleInterval = val;

        // Lower than ~100ms will crash the spark.
        if (sampleInterval < 20) {
          sampleInterval = 20;
        }
        break;

      // // Serial API
      // case msg_serialBegin:  // serial.begin
      //   type = cached[1];
      //   speed = cached[2];
      //   if (type == 0) {
      //     Serial.begin(SerialSpeed[speed]);
      //   } else {
      //     Serial1.begin(SerialSpeed[speed]);
      //   }
      //   break;

      // case msg_serialEnd:  // serial.end
      //   type = cached[1];
      //   if (type == 0) {
      //     Serial.end();
      //   } else {
      //     Serial1.end();
      //   }
      //   break;

      // case msg_serialPeek:  // serial.peek
      //   type = cached[1];
      //   if (type == 0) {
      //     val = Serial.peek();
      //   } else {
      //     val = Serial1.peek();
      //   }
      //   send(0x07, type, val);
      //   break;

      // case msg_serialAvailable:  // serial.available()
      //   type = cached[1];
      //   if (type == 0) {
      //     val = Serial.available();
      //   } else {
      //     val = Serial1.available();
      //   }
      //   send(0x07, type, val);
      //   break;

      // case msg_serialWrite:  // serial.write
      //   type = cached[1];
      //   len = cached[2];

      //   for (i = 0; i < len; i++) {
      //     if (type == 0) {
      //       Serial.write(client.read());
      //     } else {
      //       Serial1.write(client.read());
      //     }
      //   }
      //   break;

      // case msg_serialRead: // serial.read
      //   type = cached[1];
      //   if (type == 0) {
      //     val = Serial.read();
      //   } else {
      //     val = Serial1.read();
      //   }
      //   send(0x16, type, val);
      //   break;

      // case msg_serialFlush: // serial.flush
      //   type = cached[1];
      //   if (type == 0) {
      //     Serial.flush();
      //   } else {
      //     Serial1.flush();
      //   }
      //   break;

      // SPI API
      // case msg_spiBegin:  // SPI.begin
      //   SPI.begin();
      //   break;

      // case msg_spiEnd:  // SPI.end
      //   SPI.end();
      //   break;

      // case msg_spiSetBitOrder:  // SPI.setBitOrder
      //   type = cached[1];
      //   SPI.setBitOrder((type ? MSBFIRST : LSBFIRST));
      //   break;

      // case msg_spiSetClockDivider:  // SPI.setClockDivider
      //   val = cached[1];
      //   if (val == 0) {
      //     SPI.setClockDivider(SPI_CLOCK_DIV2);
      //   } else if (val == 1) {
      //     SPI.setClockDivider(SPI_CLOCK_DIV4);
      //   } else if (val == 2) {
      //     SPI.setClockDivider(SPI_CLOCK_DIV8);
      //   } else if (val == 3) {
      //     SPI.setClockDivider(SPI_CLOCK_DIV16);
      //   } else if (val == 4) {
      //     SPI.setClockDivider(SPI_CLOCK_DIV32);
      //   } else if (val == 5) {
      //     SPI.setClockDivider(SPI_CLOCK_DIV64);
      //   } else if (val == 6) {
      //     SPI.setClockDivider(SPI_CLOCK_DIV128);
      //   } else if (val == 7) {
      //     SPI.setClockDivider(SPI_CLOCK_DIV256);
      //   }
      //   break;

      // case msg_spiSetDataMode:  // SPI.setDataMode
      //   val = cached[1];
      //   if (val == 0) {
      //     SPI.setDataMode(SPI_MODE0);
      //   } else if (val == 1) {
      //     SPI.setDataMode(SPI_MODE1);
      //   } else if (val == 2) {
      //     SPI.setDataMode(SPI_MODE2);
      //   } else if (val == 3) {
      //     SPI.setDataMode(SPI_MODE3);
      //   }
      //   break;

      // case msg_spiTransfer:  // SPI.transfer
      //   val = cached[1];
      //   val = SPI.transfer(val);
      //   server.write(0x24);
      //   server.write(val);
      //   break;

      // // Wire API
      // case msg_wireBegin:  // Wire.begin
      //   address = cached[1];
      //   if (address == 0) {
      //     Wire.begin();
      //   } else {
      //     Wire.begin(address);
      //   }
      //   break;

      // case msg_wireRequestFrom:  // Wire.requestFrom
      //   address = cached[1];
      //   val = cached[2];
      //   stop = cached[3];
      //   Wire.requestFrom(address, val, stop);
      //   break;

      // case msg_wireBeginTransmission:  // Wire.beginTransmission
      //   address = cached[1];
      //   Wire.beginTransmission(address);
      //   break;

      // case msg_wireEndTransmission:  // Wire.endTransmission
      //   stop = cached[1];
      //   val = Wire.endTransmission(stop);
      //   server.write(0x33);    // could be (action)
      //   server.write(val);
      //   break;

      // case msg_wireWrite:  // Wire.write
      //   len = cached[1];
      //   uint8_t wireData[len];

      //   for (i = 0; i< len; i++) {
      //     wireData[i] = cached[1];
      //   }
      //   val = Wire.write(wireData, len);

      //   server.write(0x34);    // could be (action)
      //   server.write(val);
      //   break;

      // case msg_wireAvailable:  // Wire.available
      //   val = Wire.available();
      //   server.write(0x35);    // could be (action)
      //   server.write(val);
      //   break;

      // case msg_wireRead:  // Wire.read
      //   val = Wire.read();
      //   server.write(0x36);    // could be (action)
      //   server.write(val);
      //   break;

      case msg_servoWrite:
        pin = cached[1];
        val = cached[2];
        #ifdef DEBUG
        Serial.print("PIN: ");
        Serial.println(pin);
        Serial.print("WRITING TO SERVO: ");
        Serial.println(val);
        #endif
        servos[ToServoIndex(pin)].write(val);
        break;

      default: // noop
        break;
    } // <-- This is the end of the switch


    // Clear the cached bytes
    for (i = 0; i < bytesExpecting; i++) {
      cached[i] = 0;
    }

    // Reset hasAction flag (no longer needed for this opertion)
    // action and byte read expectation flags
    hasAction = false;
    bytesExpecting = 0;


    #ifdef DEBUG
    Serial.print("Leftovers: ");
    Serial.println(bytesRead, DEC);
    #endif

    // If there were leftover bytes available,
    // call processInput. This mechanism will continue
    // until there are no bytes available.
    if (bytesRead > 0) {
      available = bytesRead;
      processInput();
    }
  }
}


void loop() {
  if (client.connected()) {

    if (!isConnected) {
      #ifdef DEBUG
      Serial.println("--------------CONNECTED--------------");
      #endif
    }

    isConnected = true;

    // Process incoming bytes first
    available = client.available();
    if (available > 0) {

      // Move all available bytes into the buffer,
      // this avoids building up back pressure in
      // the client byte stream.
      for (int i = 0; i < available; i++) {
        buffer[i] = client.read();
        bytesRead++;
      }

      #ifdef DEBUG
      Serial.println("--------------PROCESSING NEW DATA--------------");
      #endif

      processInput();
    }

    // Reporting should be limited to every ~100ms
    nowms = millis();
    if (nowms - lastms > sampleInterval && reporters > 0) {
      lastms += sampleInterval;
      report();
    }
  } else {
    // Upon disconnection, reset the state
    if (isConnected) {
      isConnected = false;
      reset();
    }

    // If no client is yet connected, check for a new connection
    client = server.available();
  }
}
