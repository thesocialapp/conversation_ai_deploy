import unittest
from processing import transcribe_audio, generate_response  # adjust import path accordingly
from unittest.mock import patch, MagicMock

class TestProcessing(unittest.TestCase):

    @patch('processing.dg_client')  # Mocking the dg_client
    def test_transcribe_audio(self, mock_dg_client):
        # Given
        mock_audio = MagicMock()
        mock_dg_client.transcription.prerecorded.return_value = "transcript"

        # When
        result = transcribe_audio(mock_audio)

        # Then
        self.assertEqual(result, "transcript")

    def test_generate_response(self):
        # Given
        input_text = "test input"

        # When
        result = generate_response(input_text)

        # Then
        self.assertIsNotNone(result)

# to run the tests
if __name__ == '__main__':
    unittest.main()
