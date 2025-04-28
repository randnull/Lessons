import logging
import os

os.makedirs("logs", exist_ok=True)

logger = logging.getLogger('logger')
logger.setLevel(logging.DEBUG)

file_handler = logging.FileHandler("logs/notification.log", encoding='utf-8')
file_handler.setLevel(logging.DEBUG)

formatter = logging.Formatter("%(asctime)s - %(name)s - %(levelname)s - %(message)s")
file_handler.setFormatter(formatter)

if not logger.handlers:
    logger.addHandler(file_handler)