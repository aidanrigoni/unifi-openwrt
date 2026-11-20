# unifi-openwrt: OpenWrt access points as UniFi access points

Keep your old **OpenWrt** access point useful in a **UniFi** network, long after its maker stopped
supporting it. A small daemon on the AP speaks the UniFi inform protocol, so a UniFi controller adopts
the AP as a UniFi model of your choice (default: UAP-AC-Pro, `U7PG2`), pushes its settings to it, and
shows its clients, radios and stats like those of its own APs: one network, one place to manage it.

<p align="center">
  <img src="docs/images/unifi-topology.gif" alt="An OpenWrt access point in the UniFi topology view, shown as a UniFi AP" width="329">
</p>

If you like this or find it useful, support the project by buying me a coffee.

<a href="https://www.buymeacoffee.com/antoniocosta" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" alt="Buy Me A Coffee" style="height: 40px !important;width: 140px !important;" ></a>

## Why this repo exists

I own a Google OnHub (TP-Link TGR1900). It is still a capable 3x3 Wi-Fi 5 access point, but Google
ended its support in January 2023: no more updates, no more management through its app, and the
advice to recycle it. OpenWrt brings the software back up to date. My home network, though, is
UniFi, and a UniFi controller only manages Ubiquiti's own access points, so the OnHub would have
stayed a second network with its own settings, its own clients and no roaming between the two. I
don't accept that working hardware should become e-waste, or a second-class citizen in a network,
because its maker moved on.

So this repo does two things:

- **It brings OpenWrt access points into UniFi.** The AP is adopted, configured and monitored from
  the same controller as the Ubiquiti APs: Wi-Fi settings, clients, stats, roaming hints, guest
  portal, VLANs, even meshing when its cable is out. Whatever the hardware can't do is clamped and
  logged, and every change rolls back by itself if the Wi-Fi doesn't come up.
- **It documents.** How the protocol, the settings and the stats work is written up in
  [docs/internals.md](docs/internals.md), so it can be checked, ported to other devices, or carried
  on by anyone.

> **Tested on four devices only**, listed under
> [Tested devices](#tested-devices). Other OpenWrt devices, other UniFi
> models and other controller versions are **best effort**: they may work,
> they have not been tried. Reports are welcome.

**It never applies Ubiquiti firmware.** Upgrade requests from the controller
are answered and ignored; Ubiquiti firmware would brick a non-Ubiquiti
device.

## Setup: point it at your AP, pick a model, click Adopt

1. **Install the package** (OpenWrt 25.12+): build it with `tools/build.sh`
   and `tools/mkapk.sh` (see [docs/packaging.md](docs/packaging.md)), copy
   the `.apk` for your AP's architecture (`cat /etc/apk/arch`) to the AP and run
   `apk add --allow-untrusted ./unifi-openwrt_*.apk`. Then go to step 3.
   The package does not swap a basic `wpad` for the full one: it prints the
   command if you have the basic one.

   **Or** run `tools/build.sh` and copy `dist/openwrt/` (with `bin/`, `files/` and `install.sh`) to the
   AP, for example to `/tmp/uow`, and use the installer (older OpenWrt with
   opkg works too).
2. **Install** with the installer (skip this with the package):

   ```sh
   cd /tmp/uow && sh install.sh
   ```

   The installer checks free flash, replaces a basic `wpad` with the full
   one from the package feed (and restores it on failure), adds the packages
   listed under [Requirements](#requirements), installs `/usr/bin/uow`,
   `/etc/init.d/unifi-openwrt` and `/etc/config/unifi-openwrt` (an existing
   config is kept), and prints the closest UniFi models.
3. **Pick a model.** `uow -models` lists all of them; `uow -ap local
   -suggest` ranks them for this AP (same bands, then Wi-Fi generation, then
   streams). `U7PG2` is the only one with measured capability data.
4. **Configure** `/etc/config/unifi-openwrt`:

   ```sh
   uci set unifi-openwrt.main.model='U7PG2'
   uci set unifi-openwrt.main.inform_url='http://<controller IP>:8080/inform'
   uci set unifi-openwrt.main.apply='0'     # log-only at first (see Safety)
   uci set unifi-openwrt.main.enabled='1'
   uci commit unifi-openwrt
   /etc/init.d/unifi-openwrt enable
   /etc/init.d/unifi-openwrt start
   ```

   The inform URL needs an IP address, not a host name. Other options in the
   file: `hostname`, `profile_caps`, `locate_beep`, `mesh_uplink`,
   `ssh_port` (adoption SSH responder, default 2222), `hw_caps`,
   `inform_every`, `led_idle`, `uplink_port` (the port that is always the
   uplink, e.g. `wan`).
5. **Adopt.** The AP shows up in the UniFi as pending adoption, with the
   chosen model's name. Click Adopt.
6. **Watch the log:** `logread -e uow`. It shows the adoption, every plan of
   changes, what was applied, clamped or skipped, and each rollback.
   When the plans look right, set `apply '1'` and restart the service.

To stop managing the AP: Remove it in the UniFi (the daemon goes back to
pending; OpenWrt and its Wi-Fi are not reset), then
`/etc/init.d/unifi-openwrt disable; /etc/init.d/unifi-openwrt stop`.

## Tested devices

| Device | Hardware | OpenWrt | Adopted as | Tested by | What was tested |
|---|---|---|---|---|---|
| [Google OnHub (TP-Link TGR1900)](https://openwrt.org/inbox/toh/google/onhub_tp-link_tgr1900) | Qualcomm IPQ8064, Wi-Fi 5 | 25.12.6 | UAP-AC-Pro (`U7PG2`) | [@antoniocosta](https://github.com/antoniocosta) | everything under [What works](#what-works) unless marked otherwise, with UniFi Network 10.6 on a UniFi Express 7 |
| [Linksys LN1301 (MX4300)](https://openwrt.org/toh/linksys/mx4300) | Qualcomm IPQ8174 (`ath11k`), Wi-Fi 6, one 2.4 GHz and two 5 GHz radios | 25.12.6 | U6-IW | [@dohun0310](https://github.com/dohun0310) | wired, as an AP and small switch: package install and adoption, all three radios, the 160 MHz option, wired clients on the LAN ports, in-wall port numbering |
| [Linksys E4200 v2](https://openwrt.org/toh/hwdata/linksys/linksys_e4200_v2) | Marvell Kirkwood, 2x 88W8366 (`mwl8k`), Wi-Fi 4 | 23.05.6 | not reported | [@rbferreira](https://github.com/rbferreira) | package install, SSH adoption with the controller on another address, the rollback timer, mesh uplink to a UniFi AP |
| [TP-Link Archer BE800](https://forum.openwrt.org/t/openwrt-for-tp-link-be800/227055) | Qualcomm IPQ9574 (`ath12k`), Wi-Fi 7, one phy for 2.4, 5 and 6 GHz | snapshot | U7 Pro (`U7PRO`) | [@Kelon5212001](https://github.com/Kelon5212001) | adoption, all three radios, WLANs applied. Tested on 0.2.0 plus the fixes; confirmation on the merged code is pending. The mesh uplink joins but does not pass traffic (driver limit) |

Running it on another device? Open an issue or a pull request with the
device, the OpenWrt version, the UniFi model and what you tested, and it
gets a row here.

## What works

<img src="docs/images/yes.svg" alt="✓" width="16" height="16"> works &nbsp; <img src="docs/images/partial.svg" alt="△" width="16" height="16"> partial, limited or not tested live &nbsp; <img src="docs/images/no.svg" alt="✕" width="16" height="16"> not built

Tested live on the OnHub unless marked otherwise (see [Tested devices](#tested-devices)).

| UniFi feature | Status |
|---|---|
| Adoption, Remove, adopt again | <img src="docs/images/yes.svg" alt="✓" width="16" height="16"> works (SSH and inform-reply adoption) |
| Firmware update | <img src="docs/images/partial.svg" alt="△" width="16" height="16"> answered, **never applied**: nothing is downloaded; the offered version is reported so the controller stops asking (never offered in testing) |
| Clients as wireless clients of the AP, with signal, rates, channel width | <img src="docs/images/yes.svg" alt="✓" width="16" height="16"> works |
| Radio and VAP stats, channel use, interference, history graphs, CPU and memory | <img src="docs/images/yes.svg" alt="✓" width="16" height="16"> works |
| Live updates (1 s informs while the panel is open) | <img src="docs/images/yes.svg" alt="✓" width="16" height="16"> works |
| Wi-Fi experience (satisfaction) | <img src="docs/images/partial.svg" alt="△" width="16" height="16"> partial: PHY-rate and retry sub-scores, as the firmware computes them ([details](docs/internals.md#unifi-ap-satisfaction--wifi-experience-score-how-it-is-computed)) |
| Ethernet ports in the device panel | <img src="docs/images/yes.svg" alt="✓" width="16" height="16"> works (Port Manager; an AC Pro has no port squares) |
| Wired clients on the AP's LAN ports | <img src="docs/images/yes.svg" alt="✓" width="16" height="16"> works; tested by [@dohun0310](https://github.com/dohun0310) on the Linksys LN1301, not on the OnHub; in-wall models only (U6-IW and the others with the `ACCESS_POINT_SWITCH` feature) |
| Radios: channel, width, transmit power, country | <img src="docs/images/yes.svg" alt="✓" width="16" height="16"> works |
| Indoor / Outdoor | <img src="docs/images/yes.svg" alt="✓" width="16" height="16"> works (Country IE environment) |
| Minimum RSSI, load balancing (max clients) | <img src="docs/images/yes.svg" alt="✓" width="16" height="16"> works |
| WLANs: SSID, password, WPA2, WPA3, WPA2/WPA3, PMF, hidden, client isolation, U-APSD, DTIM | <img src="docs/images/yes.svg" alt="✓" width="16" height="16"> works (WPA3-only untested live) |
| WPA-Enterprise (RADIUS): WPA2, WPA3, WPA2/WPA3 | <img src="docs/images/partial.svg" alt="△" width="16" height="16"> the WLAN's first RADIUS authentication and accounting servers, PMF, PMKSA caching, RADIUS-assigned VLANs (tagged on the cable port); needs the full `wpad` |
| WPA3-Enterprise 192-bit, Enhanced Open (OWE), WEP, private pre-shared keys | <img src="docs/images/no.svg" alt="✕" width="16" height="16"> not built: such a WLAN is left off, never opened with other security |
| MAC filter (allow / deny) | <img src="docs/images/yes.svg" alt="✓" width="16" height="16"> works |
| 802.11k, 802.11v | <img src="docs/images/yes.svg" alt="✓" width="16" height="16"> works; needs the full `wpad`. The neighbour list is built from the AP's own scans: its own other-band BSSes and the same SSID on other APs of the LAN, genuine UniFi APs included |
| 802.11r (fast transition) | <img src="docs/images/partial.svg" alt="△" width="16" height="16"> switched on with OpenWrt's defaults, **untested live**. No key exchange with genuine UniFi APs, so a client moving between the two does a normal reconnect ([rrb-unifi](https://github.com/d-malinowski/unifi-openwrt-roaming) adds that exchange) |
| Band steering | <img src="docs/images/yes.svg" alt="✓" width="16" height="16"> works (usteer, local mode) |
| Minimum data rates | <img src="docs/images/partial.svg" alt="△" width="16" height="16"> clamped: applied per radio, only when every WLAN on that radio asks for it |
| WLAN schedules | <img src="docs/images/yes.svg" alt="✓" width="16" height="16"> works (enforced on the AP) |
| Rate limits (user groups) | <img src="docs/images/yes.svg" alt="✓" width="16" height="16"> works (tc shaping; nftables policing without tc) |
| Guest WLAN: isolation, guest DNS | <img src="docs/images/yes.svg" alt="✓" width="16" height="16"> works |
| Guest portal (hotspot) | <img src="docs/images/yes.svg" alt="✓" width="16" height="16"> works |
| VLANs (WLAN on a tagged network) | <img src="docs/images/yes.svg" alt="✓" width="16" height="16"> works |
| Block, Unblock, Reconnect | <img src="docs/images/yes.svg" alt="✓" width="16" height="16"> works (immediate with `hostapd-utils`) |
| Locate, LED on/off, Restart | <img src="docs/images/yes.svg" alt="✓" width="16" height="16"> works; Locate beeps on APs with a speaker |
| LED colour and brightness | <img src="docs/images/partial.svg" alt="△" width="16" height="16"> built, **untested live**: only offered when `hw_caps` is claimed (not on the AC Pro) |
| Disable device | <img src="docs/images/yes.svg" alt="✓" width="16" height="16"> works |
| IP settings (DHCP / static) | <img src="docs/images/yes.svg" alt="✓" width="16" height="16"> works; a static address other than the current one is refused |
| Device SSH authentication | <img src="docs/images/yes.svg" alt="✓" width="16" height="16"> works (extra login, root untouched) |
| NTP server, timezone | <img src="docs/images/yes.svg" alt="✓" width="16" height="16"> works |
| Remote syslog, including UniFi's encrypted syslog | <img src="docs/images/yes.svg" alt="✓" width="16" height="16"> works |
| SNMP v1/v2c | <img src="docs/images/yes.svg" alt="✓" width="16" height="16"> works; SNMPv3 not built (OpenWrt's snmpd is built without SSL) |
| RF scan, neighbour scan, rogue AP list | <img src="docs/images/yes.svg" alt="✓" width="16" height="16"> works (a spare radio is used when there is one) |
| Mesh: wireless uplink when the cable is out | <img src="docs/images/yes.svg" alt="✓" width="16" height="16"> works; never cable and mesh at once. On `mwl8k` radios (Marvell 88W8366) the uplink passes traffic only while the 5 GHz radio also carries a WLAN. On `ath12k` (reported on a TP-Link Archer BE800, OpenWrt snapshot) the uplink joins but almost no data gets through |
| Mesh Parent (choose the uplink AP) | <img src="docs/images/partial.svg" alt="△" width="16" height="16"> built, **untested live** |
| Serving mesh children, VLANs over mesh | <img src="docs/images/no.svg" alt="✕" width="16" height="16"> not built |
| Debug terminal, Support file | <img src="docs/images/no.svg" alt="✕" width="16" height="16"> not built |
| LLDP neighbours (what is on the other end of each cable) | <img src="docs/images/yes.svg" alt="✓" width="16" height="16"> works with the optional `lldpd` package, which also announces the AP to the switch: the controller then takes the AP's uplink from LLDP. lldpd must listen on the Ethernet ports themselves: `uci add_list lldpd.config.interface=lan1` (your port names), `uci commit lldpd`, `/etc/init.d/lldpd restart` |
| Discovery broadcast (the AP showing up by itself, before it has the controller's address) | <img src="docs/images/no.svg" alt="✕" width="16" height="16"> not built: the AP is given the address. A controller that finds a device this way adopts it by logging in on its SSH port 22 as `ubnt`, which this project leaves alone |
| DFS radar events | <img src="docs/images/no.svg" alt="✕" width="16" height="16"> not built |

Technical notes for contributors are in [docs/](docs/):
[docs/internals.md](docs/internals.md): protocol, config translation, stats,
the satisfaction formula, models and profiles, LED, mesh, VLANs, guest
portal and RF scan.

## Requirements

- OpenWrt with `apk` or `opkg` (tested: 25.12 and 23.05). The daemon is a static Go
  binary of 8-10 MB: APs with 16 MB of flash usually do not have room.
- **Full wpad** (`wpad-mbedtls`, `wpad-wolfssl` or `wpad-openssl`) instead of
  `wpad-basic-*`: needed for 802.11k/v and the runtime deny
  list. Install it from the package feed; check with `--simulate` before
  removing the basic one.
- `hostapd-utils` (immediate block/unblock), `kmod-nft-bridge` (guest
  isolation), `rpcd-mod-iwinfo` (`ubus call iwinfo`: without it the AP
  reports no WLANs or radio state; images with LuCI already have it).
- Optional: `usteer` (band steering), `tc-tiny`, `kmod-sched-core`,
  `kmod-sched`, `kmod-ifb` (rate-limit shaping), `snmpd` (SNMP),
  `alsa-utils` (Locate beep), `lldpd` (LLDP neighbours).
- A UniFi controller on the same network that the AP can reach on port
  8080, and that can reach the AP's adoption SSH port (default 2222).

## Safety design

- **Log-only mode**: with `apply '0'` (flag: no `-apply`) the daemon only logs
  the changes it would make.
- **Automatic rollback**: before each apply the touched UCI packages are
  backed up and an on-AP timer restores them after 5 minutes unless the
  change is confirmed. 20 s after an apply the daemon checks that every
  enabled radio with interfaces is up; if not, it rolls back at once.
- **Never wired and mesh-uplinked at once**: the mesh watcher joins only
  after 20 s without cable link and drops the wireless uplink within a
  second of a cable coming back.
- **Never applies Ubiquiti firmware**; Remove never resets OpenWrt.
- The AP never starts serving DHCP, DNS or router advertisements on its LAN:
  guest DNS runs on a separate port for guest clients only, VLAN interfaces
  have no address.
- Secrets (WLAN passphrases, the device key, password hashes) go to UCI over
  stdin, never on a command line or in the log; state files are mode 600.
- Only accounts the daemon created itself are touched for Device SSH
  authentication; root and system accounts never are.

## Building from source

Go 1.26 or later. To build just the daemon for an AP without cloning, for example an ARMv7 one:

```sh
GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 go install github.com/antoniocosta/unifi-openwrt/cmd/uow@latest
# the binary lands in $(go env GOPATH)/bin/linux_arm/uow
```

From a clone:

```sh
go test ./... && go vet ./...
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -trimpath -ldflags='-s -w' -o uow ./cmd/uow
```

`tools/build.sh` builds `uow-<arch>` for armv7, aarch64, mipsel, mips
(softfloat) and x86_64 and assembles the installer bundle in `dist/openwrt/`.

## Layout

| Path | What |
|---|---|
| `cmd/uow` | the daemon |
| `cmd/owprobe`, `cmd/owplan`, `cmd/owctl` | development helpers: print the inform fields for an AP over SSH, print the UCI plan for a saved system_cfg, run one AP action |
| `internal/` | protocol (`inform`, `device`), OpenWrt side (`openwrt`), settings translation (`translate`), model catalogue (`models`), read-only UniFi API client (`unifiapi`) |
| `packaging/` | files installed on the AP, `install.sh`, the OpenWrt package Makefile |
| `tools/` | build, `.apk` packaging, model catalogue generator, a small SNMP test client |
| `docs/` | [internals](docs/internals.md), [packaging](docs/packaging.md) |

## Contributing

Issues and pull requests are welcome. It has been tested on a few devices only (see [Tested devices](#tested-devices)), so
reports from other OpenWrt devices, UniFi models and controller versions help
most: say which device, OpenWrt version, model (`-model`) and controller version,
and attach `logread -e uow` with keys and passphrases removed (the daemon never
logs them). Run `go test ./... && go vet ./...` before sending code.

## Credits

This project builds on the work of others:

- [open-unifi](https://github.com/lucavb/open-unifi) by Luca Becker (MIT): its
  protocol documentation is the basis of the inform packet format, the
  controller replies and the system_cfg keys, and its recorded UAP-AC-Pro
  capability words are the `U7PG2` profile.
- [unifi-emu](https://github.com/jamesbraid/unifi-emu) by James Braid (MIT):
  the UniFi model catalogue (`internal/models/catalog.json`) is generated from
  its model profiles, and its protocol notes informed the inform payload.
- [unifi-inform-protocol](https://github.com/jk-5/unifi-inform-protocol) and
  [ubntmfi's inform protocol notes](https://github.com/mcrute/ubntmfi/blob/master/inform_protocol.md):
  early descriptions of the inform protocol and `set-adopt`.
- [OpenWrt](https://openwrt.org), with hostapd, iw, usteer and apk-tools, does
  the actual work on the access point.

Contributors:

- [@dohun0310](https://github.com/dohun0310): wired clients on the AP's LAN
  ports, the uplink port setting and in-wall port numbering, tested on a
  Linksys LN1301; radio reporting and the Wi-Fi 6, 802.11v and 160 MHz
  fixes.
- [@rbferreira](https://github.com/rbferreira): the rollback timer without
  `setsid`, adoption from the inform address and the mesh parent lookup,
  tested on a Linksys E4200 v2.
- [@Kelon5212001](https://github.com/Kelon5212001): radios that share one
  phy (`ath12k`) and the mesh watcher's cable check, tested on a TP-Link
  Archer BE800.

Similar projects:

- [openUF](https://github.com/yesrab/openUF) (Lua, runs on the AP; continues
  [jonasevcik/openUF](https://github.com/jonasevcik/openUF) and the original
  [trueserve/openUF](https://github.com/trueserve/openUF)): the same goal,
  with a guided installer, a LuCI app, discovery, WPA-Enterprise and more
  tested boards.
- [rrb-unifi](https://github.com/d-malinowski/unifi-openwrt-roaming): 802.11r
  fast roaming between OpenWrt and genuine UniFi access points.

Details and license texts are in [NOTICE](NOTICE).

## Legal

Not affiliated with, endorsed by or supported by Ubiquiti or the OpenWrt project. It relies on
unofficial UniFi controller behaviour that Ubiquiti can change in any update. No Ubiquiti software,
firmware or code is included, and it never installs Ubiquiti firmware. UniFi and Ubiquiti are
trademarks of Ubiquiti Inc.; OpenWrt is a trademark of its owner. The names are used only to say what
this software works with.

## License

GPL-2.0-or-later; see [LICENSE](LICENSE). Third-party attributions and license texts are in
[NOTICE](NOTICE).
