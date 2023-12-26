# Lich: The Telegram Bot for Downloading Torrents

**Disclamer**: The software author does not encourage piracy. You as the software end user are fully responsible for what you do with it.

This is a Telegram bot that runs a torrent client and manages the downloaded files. The way it works is:

 1. You drop a magnet link into the chat.
 2. The bot asks you for a category (e.g. Movies, Series, Music).
 3. The bot runs a torrent client to download the file.
 4. The bot moves the completed download into the appropritate directory based on category.
 5. You access the files (e.g. via media server such as Plex).


## Build and Install

1. Build the Debian package:

```
make lich.deb
```

2. Install the package:

```
sudo dpkg -i lich.deb
```

3. If the command above fails due to `aria2` not being installed, fix that by running:

```
sudo apt-get -f install
```

4. Edit the config file at `/opt/lich/config.json`. \
The most important fields are `token` (your bot's Telegram API token retrieved from `@BotFather`) and `whitelist` (Telegram usernames of users allowed to access the bot).

5. Enable the service to run on boot and start it:

```
sudo systemctl enable lich
sudo systemctl start lich
```

Done!
