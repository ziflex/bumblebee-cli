# bumblebee-cli
> Bumblebee app registry

CLI tool that helps to manage Bumblebee applications

## Installation

```sh
    git clone https://github.com/ziflex/bumblebee-cli
    cd ./bumblebee-cli
    make install
```

## Usage

### Add app

Register application and add prefix to dedicated ``.desktop`` file

```sh
bumblebee-cli add atom
````

It is possible to pass as many application names as needed

```sh
bumblebee-cli add atom gogland telegram slack
````

### Remove app

Unregister application and remove prefix from dedicated ``.desktop``

```sh
bumblebee-cli remove atom
````

It is possible to pass as many application names as needed

```sh
bumblebee-cli remove atom gogland telegram slack
````

### Show registered apps

Check what apps are registered and whether they are synced with their dedicated ``.desktop`` files

```sh
bumblebee-cli ls
```

See all apps in system and whether they are registered

```sh
bumblebee-cli ls -a
```

### Syncing

Update files of registered apps

```sh
bumblebee-cli sync
```

### Settings

Use another prefix

```sh
bumblebee-cli setting set prefix optirun
```