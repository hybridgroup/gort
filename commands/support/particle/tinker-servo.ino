/*
   Tinker library PLUS servo lib exposure.

   Wiring interface by Particle.io(???)

   Servo library interface by: Scott Beasley - 2016
   ---------------------------------------------------------------------

   GNU General Public License

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

// This #include statement was automatically added by the Particle IDE.
#include "Tinker-Servo.h"


/* This function is called once at start up ----------------------------------*/
void setup ( )
{
    // Init the libary, setting up function links and globals.
    init_tinker_servo ( );
}

/* This function loops forever --------------------------------------------*/
void loop ( )
{
    // All funtions are done via REST calls. Nothing needed here.
}
