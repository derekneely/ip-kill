```
._____________           ____  __.__.__  .__   
|   \______   \         |    |/ _|__|  | |  |  
|   ||     ___/  ______ |      < |  |  | |  |  
|   ||    |     /_____/ |    |  \|  |  |_|  |__
|___||____|             |____|__ \__|____/____/
                                \/             
```
  
The purpose of this application is to monitor network connectivity and changes in the public IP address. Once it notices the IP has changed or experiences a certain number of network errors it will shutdown the application passed to it. 

Note: This program has currently only been tested on an M1 Mac, Intel Mac, and Ubuntu 22.04 while running Transmission.


### How to build:
Clone the repo.
```
GOOS=[OS] GOARCH=[ARCH] go build -o /path/to/out/ip-kill
```

### How to run:

```
./ip-kill [-interval=5] [-timeout=5] [-neterr=3] process_to_kill
```

`-interval`: How often to make the external API request for the public IP. Default: 5

`-timeout`: How many seconds to allow for the API request before erroring. Default: 5

`-neterr`: The number of concurrent errors allowed before terminating. Default 3

`process_to_kill`: The process to kill if the IP changes or concurrent errors are reached.

You can also run `./ip-kill -h/-help` for a list of options.

How do you determine the 'process_to_kill' and run ip-kill? Below are a few examples.

#### M1 Mac
```
%> ps x | grep Transmission
50365   ??  S      0:00.96 /Applications/Transmission.app/Contents/MacOS/Transmission
51458 s018  S+     0:00.00 grep Transmission

%> ./ip-kill /Applications/Transmission.app/Contents/MacOS/Transmission
```

#### Ubuntu
```
$> ps x | grep transmission
 384094 ?        Sl     0:01 transmission-gtk
 384143 pts/0    S+     0:00 grep transmission

$> ./ip-kill transmission-gtk
```