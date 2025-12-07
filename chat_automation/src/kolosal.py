import os

import requests


completions_url = "https://api.kolosal.ai/v1/chat/completions"

def completions(user_prompt, system_prompt = None, max_tokens = None) -> str | None :
    # Construct required header
    headers = {
        "Authorization": f"Bearer {os.getenv("KOLOSAL_API_KEY")}",
        "Content-Type": "application/json"
    }

    # Construct request payload
    payload = {
        "model": "global.anthropic.claude-sonnet-4-5-20250929-v1:0",
        "messages": [],
    }
    if max_tokens is not None:
        payload["max_tokens"] = max_tokens
    if system_prompt is not None:
        payload["messages"].append({
            "role": "system",
            "content": system_prompt
        })
    payload["messages"].append({
        "role": "user",
        "content": user_prompt
    })

        # Send the request post
    response = requests.post(completions_url, headers=headers, json=payload)
    response.raise_for_status()

    # Return None in case error
    if response.status_code != 200:
        return None
    else:
        data_dict = response.json()
        return data_dict["choices"][0]["message"]["content"]