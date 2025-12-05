import json
import os
from dotenv import load_dotenv
from flask import Flask, request, jsonify
import requests

from chat_automation.src.kolosal import completions

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

@app.route('/webhook/<uid>', methods=['POST'])
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
        else:
            url = f"https://api.telegram.org/bot{account_api}/sendMessage"

            response = completions(user_prompt=text)

            response_payload = {
                "chat_id": chat_id,
                "text": response
            }

            response = requests.post(url, json=response_payload)
            response.raise_for_status()

            print(f"Message sent with status code: {response.status_code}")
            print("Response JSON:")
            print(json.dumps(response.json(), indent=4))

    # Return status ok & status code 200 so telegram knows we have received this.
    return jsonify({"status": "ok"}), 200



if __name__ == '__main__':
    app.run(debug=True, port=8443)