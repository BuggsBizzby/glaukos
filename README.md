<img src="https://github.com/BuggsBizzby/glaukos/blob/main/assets/glaukos.jpg" >
<h1>Glaukos: Lord of the Phishermen</h1>

Glaukos is designed to empower security professionals, educators, and enthusiasts with the capability to set up multiple isolated web interaction environments simultaneously. By harnessing the VNC-enabled Chromium service combined with mitmproxy, you can remotely engage with web content and, in parallel, observe the intricate details of user interactions with the associated websites.

<h2>Installation</h2>
<h3>Download the latest release</h3>
Move the binary to a directory with read/write permissions. Execution of the binary will create its own directory `Glaukos` to work out of during operation.

<h3>Compile your own binary</h3>
It is recommended to build the glaukos binary in a separate directory with read/write privileges, to avoid causing additional clutter in the source directory.

```
git clone [github_path]
cd glaukos
go build -o /path/to/new/directory/glaukos . 
./glaukos 
```

<h2>Usage</h2>
<h3>Summon</h3>
Run the <b>summon</b> command first to download the necessary docker images that the environments will be built off of.

```
glaukos summon
```

<h3>Create</h3>
Run the <b>create</b> command to create as many environments as your system resources will comfortably allow.
<li>-e, --environments: (Required) Specifies the number of environments to create.</li>
<li>-p, --prefix: Provides a prefix for naming environments. When provided, it's appended with a number incrementally (e.g., --prefix test will create environments named test1, test2, etc.). These names are then used as subdomains for routing purposes.</li>
<li>-n, --names: Instead of using a prefix, you can specify exact names for your environments with this comma-separated list. (e.g., intranet,sharepoint,chocolate)</li>
<li>-u, --targetURL: (Required) URL of the target website that will be displayed to the victim within the chromium service.</li>
<li>-a, --siteAddress: (Required) The domain used for routing in the Caddy config.</li>

```
glaukos create -e 3 -p document -u https://www.linkedin.com/uas/login -a yourdomain.com

glaukos create -e 3 -n intranet,sharepoint,client -u https://login.microsoftonline.com -a yourdomain.com
```

Once the environments are created you will be able to access them by navigating to them in a browser. For example, in the above context it would be: document1.yourdomain.com OR intranet.yourdomain.com

<h3>Destroy</h3>
Run the <b>destroy</b> command to tear down each of environments, either individually or all at once. Additionally, the --burn-it-down flag can be applied to remove all associated environment directories and files as well.
<li>all : Positional argument; invokes the removal of all live environments</li>
<li>-l, --list: Used to list all available environments.</li>
<li>-b, --burn-it-down: If set, environment directories and associated files will be removed.</li>

```
glaukos destroy all

glaukos destroy all -b

glaukos destroy document1

glaukos destroy document1 -b
```
