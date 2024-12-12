<div align="center">
  <a href="https://github.com/vm75/easy-share">
    <span style="font-size:40px;">Easy Share</span>
  </a>
</div>
<div align="center">

[![License]](LICENSE)
[![Build]][build_url]
[![Version]][tag_url]
[![Size]][tag_url]
[![Pulls]][hub_url]
[![Package]][pkg_url]

</div>


**Easy Share** is an open-source containerized solution for sharing file via Samba or NFS.

<!-- <p align="center">
  <img src="https://raw.githubusercontent.com/vm75/easy-share/main/docs/screenshots.gif"/>
</p> -->

## Features

- **TBD**: TBD.

## Usage  🐳

Via Docker Compose:
```yaml
services:
  easy-share:
    image: vm75/easy-share
    container_name: easy-share
    cap_add:
      - NET_ADMIN
    devices:
      - /dev/net/tun
    ports:
      - "8080:80"   # Web UI
      - 137:137
      - 138:138
      - 139:139
      - 445:445
      - 2049:2049
    volumes:
      - /path/to/data:/data
    restart: unless-stopped
```

Via Docker CLI:
```bash
docker pull vm75/easy-share
docker run -d --name easy-share \
  --cap-add=NET_ADMIN \
  --device=/dev/net/tun \
  -v /path/to/data:/data \
  -p 8080:80 \
  vm75/easy-share
```

## Configuration ⚙️


## Volume Structure

The `/data` volume should contain the following structure:
```plaintext
/data
├── config/         # Contains the sqlite3 database
├── var/            # Contains the runtime configuration and logs
```
It is recommended to place the `apps.sh` script in the `/data` volume.

## Web UI Access

The web UI is accessible at `http://<host-ip>:8080` by default. Use it to configure your samba and nfs shares with ease.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

### 3rd-Party Components

<table>
  <tr>
    <th>Component</th>
    <th>License</th>
  </tr>
  <tr>
    <td>
      <a href="https://www.samba.org/">Samba</a>
    </td>
    <td>
      <a href="https://raw.githubusercontent.com/vm75/easy-share/main/3rd-party/samba/COPYING.txt">COPYING.txt</a>
    </td>
  </tr>
  <tr>
    <td>
      <a href="https://linux-nfs.org/">nfs-utils</a>
    </td>
    <td>
      <a href="https://raw.githubusercontent.com/vm75/easy-share/main/3rd-party/nfs-utils/COPYING">COPYING</a>
    </td>
  </tr>
</table>

---

**Easy Share** provides a simple and flexible way to manage samba and nfs shares using containerization. Contributions are welcome! 🚀

[license_url]: https://github.com/vm75/easy-share/blob/main/LICENSE
[build_url]: https://github.com/vm75/easy-share/actions
[hub_url]: https://hub.docker.com/r/vm75/easy-share
[tag_url]: https://hub.docker.com/r/vm75/easy-share/tags
[pkg_url]: https://github.com/vm75/easy-share/pkgs/container/easy-share
[screenshot_url]: https://raw.githubusercontent.com/vm75/easy-share/main/docs/screenshot.gif

[License]: https://img.shields.io/badge/license-MIT-blue.svg
[Build]: https://img.shields.io/github/actions/workflow/status/vm75/easy-share/.github/workflows/ci.yml?branch=main
[Version]: https://img.shields.io/docker/v/vm75/easy-share/latest?arch=amd64&sort=semver&color=066da5
[Size]: https://img.shields.io/docker/image-size/vm75/easy-share/latest?color=066da5&label=size
[Package]: https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fipitio.github.io%2Fbackage%2Fvm75%2Feasy-share%2Feasy-share.json&query=%24.downloads&logo=github&style=flat&color=066da5&label=pulls
[Pulls]: https://img.shields.io/docker/pulls/vm75/easy-share.svg?style=flat&label=pulls&logo=docker