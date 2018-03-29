# gmina
Google text to speech API with Mina cache middleware

### Python example:

```python
import requests


def send_request(text_to_speech):
    r = requests.get("<server>/tts?text=" + text_to_speech)
    r = r.json()
    return r


if __name__ == "__main__":
    import base64
    import time
    import playsound
    import os

    text = raw_input('Text> ')

    resp = send_request(text)
    if resp["status"] == "ok":
        content = resp["content"]
        content = base64.b64decode(content)

        file_name = 'voice_{}.mp3'.format(int(time.time()))
        with open(file_name, 'wb') as f:
            f.write(content)

        playsound.playsound(file_name, True)

        os.remove(file_name)
    else:
        print resp["message"]

```
