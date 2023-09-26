import logging
import json
from datetime import datetime, timedelta
from memory import GetUser, SaveUserContext, SaveHistoryItem
from textblob import TextBlob
from functools import lru_cache

# hypothetical NLP tool for intent extraction
from nlptool import extract_intents

logger = logging.getLogger(__name__)

MAX_HISTORY_DAYS = 30


class HistoryItem:

    def __init__(self, text):
        self.text = text
        self.timestamp = datetime.now()

        # Annotate text with sentiment and intents
        self.sentiment = self.detect_sentiment(text)
        self.intents = self.detect_intents(text)

    def to_json(self):
        return json.dumps(self, default=lambda o: o.__dict__)

    @staticmethod
    @lru_cache(maxsize=None)  # Unbounded cache
    def detect_sentiment(text):
        """
        Conducts sentiment analysis on the text using TextBlob.
        Returns the polarity of the text.
        Uses caching to store previously calculated sentiment values.
        """
        return TextBlob(text).sentiment.polarity

    @staticmethod
    @lru_cache(maxsize=None)  # Unbounded cache
    def detect_intents(text):
        """
        Conducts intent extraction on the text.
        Returns a list of detected intents.
        Uses caching to store previously detected intents.
        """
        return extract_intents(text)


def prune_history(user_id):
    try:
        # Fetch user
        user = GetUser(user_id)
        if user is None or user.history is None:
            logger.error(f"No user or history found with ID {user_id}")
            return None

        # Calculate cutoff
        cutoff = datetime.now() - timedelta(days=MAX_HISTORY_DAYS)

        # Filter items older than cutoff
        user.history = [h for h in user.history if h.timestamp > cutoff]

        # Save filtered history
        SaveUserContext(user_id, user.context)
        return user.history
    except Exception as e:
        logger.error(f"Error in prune_history: {e}")
        return None


def save_response(user_id, response):
    item = HistoryItem(response)

    # Save item to database
    SaveHistoryItem(item)

    # Update context
    context = get_context(user_id, response)
    SaveUserContext(user_id, context)

    logger.info(f"Saved response for user {user_id}")
