import serial.tools.list_ports
from serial import Serial


class ISerialService:
    def connect() -> any:
        pass
    
    def read() -> any:
        pass
    
class SerialService(ISerialService):
    _serial = None
    _port = None
    
    def __init__(self, port: str) -> None:
        if len(port) != 0 and port != None:
            self._port = port
        else:
            self._port = self._find_ports()
        if len(self._port) == 0:
            raise Exception("Unable to find port")
        return
    
    def connect(self) -> any:
        try:
            self._serial = Serial(port=self._port, baudrate=115200)
        except Exception as err:
            raise err
        return
    
    def _find_ports() -> any:
        all = serial.tools.list_ports.comports()
        for port, desc, _ in sorted(all):
            if "Serial" in desc:
                return port
        return ""
    
    def read(self) -> str:
        if self._serial == None:
            raise Exception("Serial has not been initialized")
        readable: str = self._serial.in_waiting()
        if readable > 0:
            result = result + self._serial.read(readable).decode('utf-8')
            start = result.replace("!", "")
            end = result.replace("#", "")
        return
    
    def write(self) -> any:
        if self._serial == None:
            raise Exception("Serial has not been initialized")
        return
    