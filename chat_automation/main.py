import json
import os
from dotenv import load_dotenv
from flask import Flask, request, jsonify
import requests

from chat_automation.src.kolosal import completions
import chat_automation.src.telegram as telegram
from chat_automation.src.utils.utils import is_complex_json

load_dotenv()

app = Flask(__name__)

accounts = {
    "DEMO_UID": os.getenv("DEMO_API"),
}


@app.route('/', methods=['GET'])
def hello_world():
    """
    Handles a simple GET request.
    """
    response = {
        "message": "All is well",
        "status": "ok",
    }
    # jsonify serializes the Python dictionary into a JSON response
    return jsonify(response), 200

@app.route('/telegram/webhook/<uid>', methods=['POST'])
def response_webhook(uid):
    # Get the JSON Payload
    update = request.get_json()

    if 'message' in update:
        chat_id = update['message']['chat']['id']
        text = update['message']['text']
        print(f"Received message from chat '{chat_id}':{text}")
        
        account_api = accounts.get(uid)
        if account_api == None:
            print(f"Webhook came from unsaved accounts of uid '{uid}'!")
            return jsonify({"status": "ok"}), 200

        url = f"https://api.telegram.org/bot{account_api}/sendMessage"

        content = completions(user_prompt=text)

        # Parse response json to determine if its a json for backend or text for user
        if is_complex_json(content):
            #TODO pass to specific handler
            pass
        else:
            # Handle failed completions api call
            if content is None:
                content = "Maaf, sistem chatbot kami sedang ada kendala. Chat kamu sudah diteruskan ke pemilik." #TODO!

            telegram.send_message(
                api_key=account_api,
                chat_id=chat_id,
                content=content
            )
        
        # Return status ok & status code 200 so telegram knows we have received this.
        return jsonify({"status": "ok"}), 200

if __name__ == '__main__':
    app.run(debug=True, port=8443)