# -*- coding = utf-8 -*-
import logging
import config
import time


class Log:

    def __init__(self, node_id, level="DEBUG"):
        self.id = node_id
        self.logger = logging.getLogger(str(self.id))
        self.logger.setLevel(level)

    def console_handle(self, level="DEBUG"):
        console_handler = logging.StreamHandler()
        console_handler.setLevel(level)
        fmt = "[%(levelname)s] %(message)s"
        console_handler.setFormatter(logging.Formatter(fmt=fmt))
        self.logger.addHandler(console_handler)

    def file_handle(self, level="DEBUG"):
        if self.id == -1:
            filename = time.asctime() + '_' + config.consensus + '_' + config.protocol + ".log"
        else:
            filename = time.asctime() + '_' + config.consensus + '_' + config.protocol + '_' + str(self.id) + ".log"
        filename = "./log/" + filename
        file_handler = logging.FileHandler(filename=filename, mode='w', encoding="utf-8")
        file_handler.setLevel(level)
        fmt = "[%(levelname)s] %(message)s"
        file_handler.setFormatter(logging.Formatter(fmt=fmt))
        self.logger.addHandler(file_handler)

    def log_info(self, msg):
        self.logger.info(msg)

    def log_error(self, msg):
        self.logger.error(msg)
