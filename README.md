### cli-tools

####  A Simple Golang CLI Tool For Troubleshooting Inside A Container

#### Features

- **HTTP GET Request:** Fetch and display the response from any HTTP/HTTPS URL.
- **Port Scanner:** Check which TCP ports are listening on inside the container.
- **Port Checker:** Test if a specific port is open on a remote host.
- **DNS Lookup:** Query DNS records using a custom or system DNS server.

#### Build & Optimization
To build a small, statically-linked binary and further compress it for minimal size, follow these steps:

1. Build a Statically-Linked Go Binary
```sh
CGO_ENABLED=0 go build -ldflags="-s -w" -o cli-tools
```

- CGO_ENABLED=0 ensures the binary is statically linked, making it portable across Linux systems.
- -ldflags="-s -w" strips debugging information and symbol tables, reducing the binary size.

2. Compress the Binary with UPX
To achieve even smaller binaries, you can compress the executable using [UPX](https://upx.github.io/):

```sh
docker run -it --rm -v $(pwd):/app ubuntu:24.04 bash
```

```sh
apt update
apt install upx
upx --best cli-tools
exit
```
- upx --best cli-tools compresses your binary using the highest compression level.

This process ensures your Go CLI tool is as small and portable as possible, which is ideal for distribution and use in minimal environments or containers.

#### Example Usage of the Binary for Troubleshooting in a Running Kubernetes Pod

1. Copy the Binary into the Pod
First, build your binary as described earlier and ensure it is accessible on your local machine.

Then, use `kubectl cp` to copy it into the pod:

```sh
kubectl cp ./cli-tools <namespace>/<pod-name>:/tmp/cli-tools
```
- Replace \<namespace> with your pod's namespace (use default if not set).
- Replace \<pod-name> with the actual name of your pod.

Example:

```sh
kubectl cp ./cli-tools default/my-app-7c8d4d6b7f-abcde:/tmp/cli-tools
```

2. Exec Into the Pod and Run the Tool
Open a shell inside your pod:

```sh
kubectl exec -it <pod-name> -n <namespace> -- /bin/sh
```

or, if the pod uses bash:

```sh
kubectl exec -it <pod-name> -n <namespace> -- /bin/bash
```

3. Run the Binary
Inside the pod shell, make the binary executable and run it:

```sh
chmod +x /tmp/cli-tools
/tmp/cli-tools <command> [args...]
```

Example:

```sh
/tmp/cli-tools dns-lookup example.com
/tmp/cli-tools check-port google.com 443
```

**Notes:**
- This method is great for ad-hoc troubleshooting without rebuilding or redeploying your pod, especially when your pod does not have access to the Internet.
- Your binary must be statically compiled (as described above) so it doesn't depend on libraries that might not exist in the container.
- You need appropriate permissions to exec into the pod and copy files.
