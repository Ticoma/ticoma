# ![Ticoma logo](/client/assets/logo/ticoma-logo-64.png) Ticoma

## Installation guide

### Requirements

#### General

- Golang <= 1.20

#### Ubuntu

##### X11

    apt-get install libgl1-mesa-dev libxi-dev libxcursor-dev libxrandr-dev libxinerama-dev 

##### Wayland

    apt-get install libgl1-mesa-dev libwayland-dev libxkbcommon-dev

#### macOS

Get Xcode or Command Line Tools for Xcode

#### Windows

1. Download standalone gcc from [here](https://winlibs.com/)
2. Unzip the archive with [7zip](https://www.7-zip.org/)
3. Place it somewhere safe: home dir, dev folder, etc. 
4. Add the `mingw64\bin` dir to PATH

### Installation

Clone this repository

    git clone https://github.com/Ticoma/ticoma

Install raylib

    cd ticoma/main

    go get -v -u github.com/gen2brain/raylib-go/raylib

### Flags

* `client (bool)` - Run with client or headless. Default = `false`
* `relay (bool)` - When set to `true` creates a relay node instead of a player node. Default = `false`
* `fullscreen (bool)` - Run in fullscreen mode instead of a quarter resolution window. Default = `false`

### .Env file

To successfully connect to the libp2p pubsub and join the game,  
you will need a proper .env file in the `main` folder.

Currently the .env file consists of:

* `DEBUG` - Debug mode.
```
0 = NO LOGS
1 = ALL LOGS
2 = ONLY LOGS RELATED TO PLAYER
3 = ONLY LOGS RELATED TO NETWORK
```
* `RELAY_ADDR` - PeerID of Relay node to join pubsub
* `RELAY_IP` - Public IPv4 address of relay node