from dataclasses import dataclass
from datetime import datetime, timedelta
import threading
import time

# --- Configuration ---
TIMEOUT_SECONDS = 60 * 60 * 4 # TTL for message history used as context for the AI. 
MESSAGE_LEN_LIMIT = 20 # Limit of messages stored for a single chat room

# --- Data struct ---
@dataclass
class Message:
    """
    A single message entity

    Attributes:
        sender_id (int): Message sender ID (user_id).
        content (string): Content of the message.
        timestamp (string): Timestamp of the message stored in ISO Format.
    """
    sender_id: int # Either user_id or 0 for bot
    content: str
    timestamp: str # Stored as ISO format for clearance when sending to AI

@dataclass
class Chat:
    """
    A single chat entity to store message history. To be used on a dictionary keyed by its chat_id

    Attributes:
        message (list[Message]): message history list
        last_active (datetime): chat's last active based on latest message sent to the chat room. 
    """
    messages: list[Message]
    last_active: datetime

# --- Store ---
lock = threading.Lock()
message_history: dict[int, Chat] = {}

_CLEANUP_THREAD_STARTED = False 

def add_message(chat_id: int, sender: int, content: str):
    """
    Add a message to the history
    
    :param chat_id: Telegram chat room ID
    :type chat_id: int
    :param sender: Sender ID. For simplification use 0 to declare BOT response message
    :type sender: int
    :param content: Content of the message
    :type content: str
    """
    current_time = datetime.now()

    new_message = Message(
        sender_id=sender,
        content=content,
        timestamp=current_time.isoformat()
    )

    with lock:
        # Initialize chat room if it doesn't exist
        if chat_id not in message_history:
            message_history[chat_id] = Chat(
                messages = [],
                last_active=current_time
            )

        chat = message_history[chat_id]

        # 1. Add message and reset last active time
        chat.messages.append(new_message)
        chat.last_active = current_time

        # 2. Enforce MESSAGE_LEN_LIMIT (keeping only the latest messages)
        if len(chat.messages) > MESSAGE_LEN_LIMIT:
            chat.messages = chat.messages[-MESSAGE_LEN_LIMIT:]
            
        print(f"[LOG] Added message for chat room: {chat_id}. History size: {len(chat.messages)}")

def get_history(chat_id: int) -> list[Message]:
    """
    Docstring for get_history
    
    :param chat_id: Telegram chat room ID
    :type chat_id: int
    :return: List of message history
    :rtype: list[Message]
    """
    with lock:
        if chat_id not in message_history:
            print(f"[LOG] No history found for {chat_id}.")
            return []
        
        # Reset the inactivity timer when history is accessed
        chat = message_history[chat_id]
        chat.last_active = datetime.now()
        print(f"[LOG] Retrieved history for {chat_id} and reset inactivity timer.")

        return chat.messages

def delete_history(chat_id: int):
    """
    Delete message history from a chat room.
    
    :param chat_id: Telegram chat room ID
    :type chat_id: int
    """
    with lock:
        if chat_id in message_history:
            del message_history[chat_id]
            print(f"[LOG] Manually deleted history for {chat_id}.")
        else:
            print(f"[LOG] Cannot delete history for {chat_id} because it doesn't exist.")

def cleanup_history_task():
    """
    A persistent background task that deletes inactive user histories.
    """
    inactivity_threshold = timedelta(seconds=TIMEOUT_SECONDS)
    
    # Create a sensible sleep interval
    sleep_interval = 60 * 30
    
    while True:
        try:
            with lock:
                current_time = datetime.now()

                # Identify users to delete
                chats_to_delete = []
                for chat_id, chat_data in message_history.items():
                    # FIX: Access .last_active attribute on the Chat object
                    if current_time - chat_data.last_active > inactivity_threshold:
                         chats_to_delete.append(chat_id)
                
                # Delete them
                for chat_id in chats_to_delete:
                    del message_history[chat_id]
                    print(f"[CLEANUP] Cleaned up inactive history for {chat_id}.")   
            
            # Print total active chats outside the lock (or use a dedicated logger)
            print(f"[CLEANUP] Total active chats remaining: {len(message_history)}")

        except RuntimeError as e:
            # Handle dictionary size change error (should be mitigated by the lock, but safer to catch)
            print(f"[CLEANUP ERROR] Handled dict iteration error: {e}")
            pass

        # Sleep for the defined interval
        time.sleep(sleep_interval) 

def start_history_cleanup_service():
    """
    Initializes and starts the background cleanup thread once.
    This function should be called by the main application entry point (e.g., main.py).
    """
    global _CLEANUP_THREAD_STARTED
    if _CLEANUP_THREAD_STARTED:
        return

    # Start the background cleanup thread
    # daemon=True ensures the thread is terminated when the main program exits
    cleanup_thread = threading.Thread(target=cleanup_history_task, daemon=True)
    cleanup_thread.start()
    _CLEANUP_THREAD_STARTED = True
    print(f"History cleanup service started. TTL is {TIMEOUT_SECONDS} seconds ({TIMEOUT_SECONDS / 60 / 60} hours).")
