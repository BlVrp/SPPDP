# SPPDP

How to run frontend:

1. cd one-help-fe
2. npm install
3. npx expo start

How to run on your phone( if you have Android):

1. Enable developer options:

Go to Settings > About phone (or System > About phone).
Tap Build number (or Software version) 7 times until you see a message saying “You are now a developer!”

2. Enable USB debugging:

Go to Settings > Developer options.
Turn on USB debugging.

3. Connect your phone to your computer via USB cable.

Make sure your phone is unlocked, and allow any prompts to trust your computer.

## How to run backend...

copy config file content of `.one-help.env` into `./backend/configs/.one-help.env`

then run
```bash
cd ./backend ; docker compose up --build
```
