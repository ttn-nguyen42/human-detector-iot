import serial.tools.list_ports
from serial import Serial


class ISerialService:
    def connect(self) -> None:
        pass

    def read(self) -> str:
        pass

    def write(self, payload: any) -> None:
        pass


class SerialService(ISerialService):
    _serial = None
    _port = None

    def __init__(self, port: str) -> None:
        if len(port) is not 0 and port is not None:
            self._port = port
        else:
            self._port = self._find_ports()
        if len(self._port) == 0:
            raise Exception("Unable to find port")
        return

    def connect(self) -> None:
        try:
            self._serial = Serial(port=self._port, baudrate=115200)
        except Exception as err:
            raise err
        return

    def _find_ports(self) -> str:
        all = serial.tools.list_ports.comports()
        for port, desc, _ in sorted(all):
            if "Serial" in desc:
                return port
        return ""

    def read(self) -> str:
        if self._serial is None:
            raise Exception("Serial has not been initialized")
        readable: str = self._serial.in_waiting()
        if readable > 0:
            result: str = self._serial.read(readable).decode('utf-8')
            start = result.replace("!", "")
            end = result.replace("#", "")
            # Gets clean result string here
        return result

    def write(self, payload: any) -> None:
        if self._serial is None:
            raise Exception("Serial has not been initialized")
        try:
            self.write((str(payload) + "#").encode())
        except Exception as err:
            raise err
        return
