import io
from typing import Union
import requests


def send_message(
    api_key: str,
    chat_id: Union[int, str],
    content: str
) -> bool:
    """
    Sends a message to a Telegram chat room.

    :param api_key: The API Key of the Telegram Bot.
    :param chat_id: The ID of the target chat.
    :param content: The message content being sent.
    :return: True if the request was successful, False otherwise.
    """
    url = f"https://api.telegram.org/bot{api_key}/sendMessage"
    message_payload = {
        "chat_id": chat_id,
        "text": content
    }

    try:
        print(f"\n[TELEGRAM] Attempting to send message to chat {chat_id}...")

        response = requests.post(url, json=message_payload)
        response.raise_for_status()

        if response.status_code == 200:
            print(f"Telegram Success: Message sent.")
            return True
        else:
            error_desc = response.json().get('description', 'Unknown error.')
            print(f"Telegram Error: Status {response.status_code}. Details: {error_desc}")
            print("NOTE: Did you replace 'YOUR_BOT_TOKEN' and ensure the chat ID is correct?")
            return False
    except requests.exceptions.RequestException as e:
        print(f"Network Error: Could not connect to Telegram API. Error: {e}")
        return False
    except Exception as e:
        print(f"An unexpected error occurred: {e}")
        return False


def send_photo(
    api_key: str,
    chat_id: Union[int, str], 
    photo_source: Union[str, io.BytesIO], 
    caption: str = None
) -> bool:
    """
    Sends a photo to a Telegram chat room.

    Supports three methods for the photo_source:
    1. String (HTTP URL or file_id): Passes the string directly in the request payload.
    2. io.BytesIO: Uploads the in-memory image using multipart/form-data.

    :param api_key: The API Key of the Telegram Bot.
    :param chat_id: The ID of the target chat.
    :param photo_source: The source of the photo (URL, file_id, or BytesIO object).
    :param caption: Optional photo caption.
    :return: True if the request was successful, False otherwise.
    """
    url = f"https://api.telegram.org/bot{api_key}/sendPhoto"
    payload = {'chat_id': chat_id}
    files = {}

    if caption:
        payload['caption'] = caption

    if isinstance(photo_source, str):
        # 1. HTTP URL or Telegram File ID
        payload['photo'] = photo_source
        
    elif isinstance(photo_source, io.BytesIO):
        # 2. Multipart/form-data upload (in-memory file)
        # 'photo' is the field name expected by Telegram's sendPhoto method
        files = {'photo': ('qr_code.png', photo_source, 'image/png')}
    
    try:
        print(f"\n[TELEGRAM] Attempting to send photo to chat {chat_id}...")
        
        response = requests.post(
            url, 
            data=payload, 
            files=files if files else None # Only include 'files' if uploading
        )
        
        # Check the response status
        response_json = response.json()
        if response.status_code == 200:
            print(f"Telegram Success: Photo sent.")
            return True
        else:
            error_desc = response_json.get('description', 'Unknown error.')
            print(f"Telegram Error: Status {response.status_code}. Details: {error_desc}")
            return False

    except requests.exceptions.RequestException as e:
        print(f"Network Error: Could not connect to Telegram API. Error: {e}")
        return False
    except Exception as e:
        print(f"An unexpected error occurred: {e}")
        return False