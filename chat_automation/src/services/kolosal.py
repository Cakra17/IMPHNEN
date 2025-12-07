import os

import requests

class KolosalService:
    token: str
    completions_url: str

    def __init__(self, token):
        self.token = token
        self.completions_url = "https://api.kolosal.ai/v1/chat/completions"

    def completions(self, user_prompt, system_prompt = None, max_tokens = None) -> str | None :
        try:
            headers = {
                "Authorization": f"Bearer {self.token}",
                "Content-Type": "application/json"
            }

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

            response = requests.post(self.completions_url, headers=headers, json=payload, timeout=30)
            response.raise_for_status()

            if response.status_code != 200:
                return None
            else:
                data_dict = response.json()
                return data_dict["choices"][0]["message"]["content"]
        except requests.exceptions.Timeout:
            return None
        except requests.exceptions.ConnectionError:
            return None
        except requests.exceptions.HTTPError:
            return None
        except Exception:
            return None