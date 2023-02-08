# Automatic testing

## Prepare the test environment

First prepare your environment. Set up a freeipa server either in a dedicated host or in a container. 
For container, you can refer to the [Test Workflow](.github/workflows/test-acc.yml) for some insight.

The freeipa server must be resolvable either by dns or using `/etc/hosts`.

## Run the tests

First, set up the environment variables:
```bash
export FREEIPA_HOST="ipa.ipatest.lan"
export FREEIPA_USERNAME="admin"
export FREEIPA_PASSWORD="P@ssword123"
```
Then run the test:
```bash
TF_ACC=1 go test -v -cover ./...
```

## Run FreeIPA server in podman

Allow cgroup in the container:
```bash
sudo setsebool -P container_manage_cgroup 1
```
Allow to bind on privileged port > 80 
```bash
sudo sysctl -w net.ipv4.ip_unprivileged_port_start=80
```

Run FreeIPA server in podman
```bash
export FREEIPA_VERSION="fedora-37-4.10.1"
podman container run -ti --rm --name freeipa-server-container -h ipa.ipatest.lan --dns=127.0.0.1 --read-only -p 127.0.0.1:80:80 -p 127.0.0.1:443:443 -v $HOME/Tmp/ipa-data:/data:Z -e container=podman -e IPA_SERVER_HOSTNAME="ipa.ipatest.lan" -e IPA_SERVER_INSTALL_OPTS='--no-ntp --ds-password=P@ssword123 --admin-password=P@ssword123 --domain=ipatest.lan --realm=IPATEST.LAN --no-host-dns --no-forwarders --setup-dns --no-dnssec-validation --allow-zone-overlap --no-reverse --unattended' freeipa/freeipa-server:${FREEIPA_VERSION}
```

Once the container initialized, you can update `/etc/hosts` to add the ipa hostname for ip 127.0.0.1. If done before, the container initialization will fail.

`/etc/hosts` needs to be updated before and after the container start otherwise it will not initialize correctly.

Setting the hostname in your external DNS server works however.
