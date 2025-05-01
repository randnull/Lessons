from datetime import datetime

from .job_config import scheduler

from AnswerEngine.src.config.settings import settings
from AnswerEngine.src.logger.logger import logger
from .job_functions import selected_status_check


async def start_scheduler():
    delay: int = settings.DELAY_SCHEDULE

    logger.info(f"Job started at {datetime.now()}")
    scheduler.add_job(selected_status_check, 'interval', seconds=delay)
    scheduler.start()


async def stop_scheduler():
    logger.info(f"Job stopped at {datetime.now()}")
    scheduler.shutdown()
