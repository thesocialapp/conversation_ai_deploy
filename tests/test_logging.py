import io
import json
import logging
from myproject.main import AppLoggingFormatter  # adjust the import according to your project structure


def test_app_logging_formatter():
    logger = logging.getLogger('test_logger')
    log_output = io.StringIO()
    handler = logging.StreamHandler(log_output)
    handler.setFormatter(AppLoggingFormatter())
    logger.addHandler(handler)
    logger.error("This is a test error message", extra={'trace_id': '1234'})
    logged_output = json.loads(log_output.getvalue())
    assert logged_output['level'] == 'ERROR'
    assert logged_output['message'] == "This is a test error message"
    assert logged_output['trace_id'] == '1234'
