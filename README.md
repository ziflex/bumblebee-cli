# bumblebee-gnome
CLI tool that helps to manage Bumblebee dependant applications

# Usage

## Add app

Register application and add prefix to dedicated ``.desktop`` file

```sh
bumblebee-gnome add atom
````

It is possible to pass as many application names as needed

```sh
bumblebee-gnome add atom gogland telegram slack
````

## Remove app

Unregister application and remove prefix from dedicated ``.desktop``

```sh
bumblebee-gnome remove atom
````

It is possible to pass as many application names as needed

```sh
bumblebee-gnome remove atom gogland telegram slack
````

## Show registered apps

Check what apps are registered and whether they are synced with their dedicated ``.desktop`` files

```sh
bumblebee-gnome ls
```

See all apps in system and whether they are registered

```sh
bumblebee-gnome ls -a
```

## Syncing

Update files of registered apps

```sh
bumblebee-gnome sync
```
