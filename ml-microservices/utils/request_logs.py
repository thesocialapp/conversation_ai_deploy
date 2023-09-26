import logging
from time import perf_counter
import uuid  # import uuid for generating request IDs

logger = logging.getLogger(__name__)

class RequestLogger:
    def __init__(self, app):
        self.app = app

    def __call__(self, environ, start_response):
        request_data = {
            # ... populate with request metadata
        }

        request_id = str(uuid.uuid4())  # Generate request IDs
        request_data["request_id"] = request_id
        logger.info(f"Request: {request_data}")  # Use f-strings for logging values

        start = perf_counter()

        def logging_start_response(status, headers, *args):
            request_data["status"] = status
            request_data["response_time"] = perf_counter() - start
            logger.info(f"Response: {request_data}")  # Use f-strings for logging values
            return start_response(status, headers, *args)

        return self.app(environ, logging_start_response)
