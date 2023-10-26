<img src="https://github.com/BuggsBizzby/glaukos/blob/main/assets/glaukos.jpg" >
<h1>Glaukos: Lord of the Phishermen</h1>

Glaukos is designed to empower security professionals, educators, and enthusiasts with the capability to set up multiple isolated web interaction environments simultaneously. By harnessing the VNC-enabled Chromium service combined with mitmproxy, you can remotely engage with web content and, in parallel, observe the intricate details of user interactions with the associated websites.

<h2>Installation</h2>
<h3>Download the latest release</h3>
Move the binary to a directory with read/write permissions. Execution of the binary will create its own directory `Glaukos` to work out of during operation.

<h3>Compile your own binary</h3>

```
git clone [github_path]
cd glaukos
go build -o /path/to/new/directory/glaukos . ----> This will compile the go binary, name it `glaukos` and store it in a directory you define
./glaukos ----> Execute the binary
```

<h2>Usage</h2>
