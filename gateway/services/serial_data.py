import serial.tools.list_ports
from serial import Serial
import logging


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
        if len(port) != 0 and port != None:
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

    def read(self):
        if self._serial is None:
            raise Exception("Serial has not been initialized")
        readable: int = self._serial.in_waiting
        result = []
        if readable > 0:
            try:
                raw_seq: str = self._serial.read(readable).decode('utf-8')
            except Exception as err:
                logging.info(f"Cannot read from serial, skipping receiving data err={err}")
                return result
            while ('{' in raw_seq) and ('}' in raw_seq):
                start = raw_seq.find('{')
                end = raw_seq.find('}')
                if start < 0 or end < 0:
                    return result
                result.append(raw_seq[start:end +1])
                if end == len(raw_seq):
                    raw_seq = ""
                    continue
                raw_seq = raw_seq[end + 1:]
        # Returns a JSON
        return result

    def write(self, payload: any) -> None:
        if self._serial is None:
            raise Exception("Serial has not been initialized")
        try:
            self.write((str(payload)).encode())
        except Exception as err:
            raise err
        return
