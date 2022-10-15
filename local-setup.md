# Local setup

### Why?

As of the time of writing I'm approximately 25 minutes into having local go links working, so I'm not too sure yet. At least so far it seems good for where you need the path part of a URL, but browser autocomplete only works up to domain (Github, for example). Also handy for links to documentation/blogs with obscure domains. I can also see this being handy for saving and naming links as you find them, and rediscovering them by browsing `~/.gosearch` (although now that I think of it this just sounds like a roundabout way to bookmark stuff).

### Prerequisites

Have a laptop setup identically to mine.

### More Abstract Prerequisites

* Git
* Docker
* SystemD
* Linux (probably)

## Steps

1. Clone this repo

```
git clone https://github.com/ckw017/go-shorten
```

2. Build and tag an image

```
cd go-shorten
docker build -t go-shorten .
```

3. Make a directory wherever. We'll be mounting this into the container so that we can have persistent storage.

```
mkdir ~/.gosearch
```

4. Create a systemd unit file. Remember to update `<this-you>` to your user.

```
# /etc/systemd/system/go-shorten.service

[Unit]
Description=Local go shorten service
After=network.target
StartLimitIntervalSec=0

[Service]
Type=simple
Restart=always
RestartSec=1
User=chriskw
ExecStartPre=-docker stop go-shorten-container
ExecStart=docker run --rm -v /home/<this-you>/.gosearch/:/var/gosearch -p 80:8080 go-shorten --root-path=/var/gosearch --storage-type=filesystem

[Install]
WantedBy=multi-user.target
```

5. Start and enable the service

```
sudo systemctl start go-shorten
sudo systemctl enable go-shorten
```

You should be able to see the page in browser at `localhost:80` now.

6. Add the following line to `/etc/hosts`:

```
127.0.0.1 localhost go
```

If there's already a line that looks like that but without the `go`, you can replace it.

7. Configure your browser to recognize `go` as a domain name. Without this step, your browser might think you're just trying to search for the word "go". I followed the instructions [here](https://support.mozilla.org/en-US/questions/1285922). You can probably find similar instructions by searching for stuff like "name-of-your-browser use local TLD".

After this you should be good to `go/`. It looks like search is borked when using filesystem storage, so in the meantime you can do `ls ~/.gosearch` to see what you've got.
