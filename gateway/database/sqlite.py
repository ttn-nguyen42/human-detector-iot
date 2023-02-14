import logging
import os
import sqlite3

ENV_SQLITE_DATABASE="SQLITE_DATABASE"

def get_sqlite_database() -> str:
    # Reads SQLITE_DATABASE environment variable
    res = os.environ.get(ENV_SQLITE_DATABASE)
    if len(res) == 0:
        raise Exception("SQLite database environment variable not found")
    return res

class SQLDatabase:
    # Interface for databases
    def connect(self) -> any:
        pass

    def execute(self, query: str) -> any:
        pass

    def commit(self) -> None:
        pass

    def close(self) -> None:
        pass


class SqliteDatabase(SQLDatabase):
    # Implements SQLDatabase
    # Connects to a database
    # then provide APIs for executing SQL queries

    _conn: sqlite3.Connection = None
    _database_name: str = ""

    def __init__(self, which: str) -> None:
        self._database_name = which
        return

    def connect(self) -> any:
        self._conn = sqlite3.connect(self._database_name)
        logging.info(f"SQLite database {self._database_name} connected")
        return self

    def execute(self, query: str) -> any:
        try:
            cursor = self._conn.execute(query)
        except Exception as err:
            raise err
        return cursor
    
    def commit(self) -> None:
        self._conn.commit()
        return
    
    def close(self) -> None:
        self._conn.close()
        logging.info("SQLite database connection closed")
        return
