# bumblebee-ui
CLI tool that helps to manage Bumblebee dependant applications

## Installation

```sh
    git clone https://github.com/ziflex/bumblebee-ui
    cd ./bumblebee-ui
    make install
```

## Usage

### Add app

Register application and add prefix to dedicated ``.desktop`` file

```sh
bumblebee-ui add atom
````

It is possible to pass as many application names as needed

```sh
bumblebee-ui add atom gogland telegram slack
````

### Remove app

Unregister application and remove prefix from dedicated ``.desktop``

```sh
bumblebee-ui remove atom
````

It is possible to pass as many application names as needed

```sh
bumblebee-ui remove atom gogland telegram slack
````

### Show registered apps

Check what apps are registered and whether they are synced with their dedicated ``.desktop`` files

```sh
bumblebee-ui ls
```

See all apps in system and whether they are registered

```sh
bumblebee-ui ls -a
```

### Syncing

Update files of registered apps

```sh
bumblebee-ui sync
```
