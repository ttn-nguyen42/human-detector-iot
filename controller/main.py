from yolobit import *
import uselect
import time
import sys
import uasyncio as asyncio
import dht
from machine import Pin
import json
import music
from event_manager import *
import _thread

#
# Author: nguyen_tran
#

##################################################################


class SerialService:
  #
  # Through USB serial
  #
  _spoll = None

  def __init__(self):
    self._spoll = uselect.poll()        # Set up an input polling object.
    # Register polling object.
    self._spoll.register(sys.stdin, uselect.POLLIN)
    return

  def write(self, inp: str):
    print(inp)
    return

  def read(self) -> str:
    inp = ''

    if self._spoll.poll(0):
      inp = sys.stdin.read(1)

    while self._spoll.poll(0):
      inp = inp + sys.stdin.read(1)

    return str(inp)

  def __del__(self):
    self._spoll.unregister(sys.stdin)
    return

##################################################################


class DHTNode:
  _dht11 = None

  def __init__(self, pin: Pin):
    self._dht11 = dht.DHT11(pin)
    return

  def read(self):
    # Returns temperature and humidity
    temp = self._dht11.temperature()
    humid = self._dht11.humidity()
    self._dht11.measure()
    return temp, humid


class DetectorNode:
  _status = False

  def __init__(self):
    self._status = False
    return

  def read(self):
    # Check for motion
    sig = pin0.read_analog()
    if sig > 0:
      # Found human
      # Display green circle
      display.show(Image("00000:04440:04040:04440:00000"))
      # Play once
      if self._status == False:
        music.play(['C6:2'], wait=False)
      self._status = True
    else:
      display.show(Image("00000:01110:01010:01110:00000"))
      music.stop()
      self._status = False
    return sig > 0

##################################################################


class SensorDto:
  temp: int
  humidity: int
  detected: bool

  def __init__(self, temp: int, humidity: int, detected: bool):
    self.temp = temp
    self.humidity = humidity
    self.detected = detected

##################################################################


def send_data(data: SensorDto):
  # Sends as JSON
  js = json.dumps(data.__dict__)
  print(js)
  return


def execute_command(command: str):
  if command == "A":
    music.play(['A1:2'], wait=False)
  return

##################################################################


serial_service = SerialService()

thread_must_stop = False


def read_commands():
  while True:
    # print("Checking command")
    task = serial_service.read()
    execute_command(task)
    time.sleep(0.5)
  return


def get_sensor_data(dht_service, detector_service):
  print("Sending sensor data")
  while True:
    if thread_must_stop == True:
      print("Stop sensor data")
      return
    temp, humid = dht_service.read()
    detected = detector_service.read()
    dto = SensorDto(temp, humid, detected)
    send_data(dto)
    time.sleep(1)
    continue
  return

##################################################################


def main():
  # Setup
  global thread_must_stop
  dht_node = DHTNode(Pin(pin1.pin))
  detector_node = DetectorNode()

  # Welcome
  music.play(music.POWER_UP, wait=True)
  display.show(Image("00000:01110:01010:01110:00000"))
  print("Starting")

  # Creating tasks
  _thread.start_new_thread(get_sensor_data, (dht_node, detector_node))
  try:
    read_commands()
  except KeyboardInterrupt:
    thread_must_stop = True
    sys.exit()
  return


##################################################################
main()