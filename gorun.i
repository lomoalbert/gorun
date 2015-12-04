%module gorun
%{
/* Put header files here or function declarations like below */

#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <stdint.h>
#include <string.h>
#include <math.h>
#include "I2Cdev.h"
#include "helper_3dmath.h"
#include "MPU6050_6Axis_MotionApps20.h"
#include "MPU6050.h"
%}

%include "MPU6050.h"
%include "MPU6050_6Axis_MotionApps20.h"