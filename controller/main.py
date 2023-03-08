from yolobit import *
import uselect
import time
import sys
import uasyncio as asyncio
import dht
from machine import Pin
import json

#
# Author: nguyen_tran
# This uses MicroPython as opposed to Python
#

##################################################################


class ICommunicationService:
  #
  # Generic interface for communication methods
  #

  # Writes to the communication stream
  def write(self, inp: str):
    pass

  # Read from the communication stream
  def read(self) -> str:
    pass

##################################################################


class SerialService(ICommunicationService):
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


serial_service: ICommunicationService = SerialService()

##################################################################


class ISensorNode:
  # Gets sensor data
  def read():
    pass


class DHTNode(ISensorNode):
  def __init__(self):
    return

  def read():
    dht11.measure()
    # Returns temperature and humidity
    temp = dht11.temperature()
    humid = dht11.humidity()
    return (temp, humid)


dht_node = DHTNode()


class DetectorNode(ISensorNode):
  def __init__():
    return

  def read():
    # Check for motion
    return False


detector_node = DetectorNode()

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


async def send_data(data: SensorDto):
  # Sends as JSON
  js = json.dumps(data)
  print(js)
  return


async def execute_command(command: str):
  print("Received: {0}".format(task))
  return

##################################################################


async def read_commands(service: ICommunicationService):
  while True:
    task = service.read()
    execute_command(task)
    await asyncio.sleep(500)
    continue
  return


async def get_sensor_data(dht_node: ISensorNode, detector_node: ISensorNode):
  while True:
      temp, humid = dht_node.read()
      detected = detector_node.read()
      dto = SensorDto(temp, humid, detected)
      task = asyncio.create_task(send_data(dto))
      await asyncio.sleep(500)
      continue
  return

##################################################################


async def main():
  # Create tasks
  # One for reading the commands
  command_reader = asyncio.create_task(read_commands(serial_service))
  # One for reading the sensor and send it through serial
  data_reader = asyncio.create_task(get_sensor_data(dht_node, detector_node))
  return

##################################################################

asyncio.run(main())