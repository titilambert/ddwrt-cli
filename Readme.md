# DDWRT-CLI

This is a small cli to manage quickly your dd-wrt router

## Build

```
make build
```

## Usage

### Init

You need to run the `init` command before to be able to use ddwrt-cli:
```
ddwrt-cli -H ROUTER_IP -u USERNAME -p PASSWORD init
```

Each time you want to connect to an other router, you need to rerun the init command.

### Show current config

Then you need to dump the configuration of the ddwrt web page.
For example, to dump the page `Setup`=>`Basic Setup`, you have to run:
```
ddwrt-cli H ROUTER_IP -u USERNAME -p PASSWORD setup basic-setup
```

### Dump config

Then you need to dump the configuration of the ddwrt web page.
For example, to dump the page `Services`=>`Services`, you have to run:
```
ddwrt-cli H ROUTER_IP -u USERNAME -p PASSWORD services services -d dump/router_name/
```

### Edit the config

Edit the config file and change what you want:
```
vim router_name/services/services.yaml
```

### Load config

Then you can apply your changes with the following command:
```
ddwrt-cli H ROUTER_IP -u USERNAME -p PASSWORD services services -l dump/router_name/
```


## Notes

I'm not a Golang expert, so I'm pretty sure this code can be improved.

If you want to improve it, please, submit a PR ! Thanks

Tested on the following pages:
* `Services` => `Services`

Tested on the following routers:
* Router Model: Dlink DIR-825
  Firmware Version: DD-WRT v3.0-r34411 std (01/07/18)
* Router Model: Dlink DIR-825
  Firmware Version: DD-WRT v3.0-r36079 std (06/01/18)
* Router Model: Buffalo WZR-HP-G300NH2
  Firmware Version: DD-WRT v3.0-r33679 std (11/04/17)
* Router Model: Buffalo WZR-HP-G300NH2
  Firmware Version DD-WRT v3.0-r34411 std (01/07/18)
