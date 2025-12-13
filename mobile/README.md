# Flutter Web-in-Docker (Web-First, APK-at-End)

This README explains how to run a Flutter app in a Docker container using a web workflow. The app is served on port **8080**, and you can view it at **http://localhost:8080** from your host browser.

---

## Prerequisites
- **Docker** and **Make** installed on the host.
- Project folder on host: `mobilePiscine/` (mounted into the container at `/app`).
- A Flutter app folder named `ex00` inside `mobilePiscine/` (create it on first run if needed).

---

## Makefile targets (summary)

### 1) Build the Docker image
```bash
make build
```
- Builds the image tagged as `$(NAME)` (ensure `NAME := your-image-tag` is set in the Makefile or in your shell).
- Installs and configures the Flutter/Android SDK toolchain inside the image (per your Dockerfile).

### 2) Start a generic shell in the container
```bash
make run
```
- Opens an interactive shell (`bash`) in the container.
- Mounts your host folder `$(PWD)/mobilePiscine` to `/app` inside the container.
- Port `8080` is published for you to run a dev server manually if you want.

### 3) Start the web dev server (preferred for web workflow)
```bash
make web
```
- Runs `flutter config --enable-web` (idempotent).
- Changes into `ex00` and starts:
  ```bash
  flutter run -d web-server --web-port 8080 --web-hostname 0.0.0.0
  ```
- The app is now served from the container on **port 8080**.

---

## View the app
Open your host browser at:

**http://localhost:8080**

You should see the default Flutter template (or your app UI if you’ve already modified it).

---

## First-time setup (if `ex00` doesn’t exist yet)

If the `ex00` folder is missing in `mobilePiscine/`, you can create it:

**Option A — via `make run`:**
```bash
make run
# then inside the container:
cd /app
flutter config --enable-web
flutter create ex00
cd ex00
flutter run -d web-server --web-port 8080 --web-hostname 0.0.0.0
```

**Option B — adjust your Makefile** to create `ex00` automatically before running.

---

## Stopping the server
In the terminal that’s running `flutter run`:
- Press `q` to stop, or
- `Ctrl + C` to interrupt the process.

The container started by `make web` will exit when the command ends.

---

## Hot reload / Hot restart (current status: **not working**)
At the moment, pressing `r` (hot reload) or `R` (hot restart) in the `flutter run` terminal **times out**. This is **unresolved** in the current setup.

**Why it may be failing (brief):**
- The `web-server` device runs a **Dart VM service** bound to `127.0.0.1` **inside the container**; your **host browser** can’t attach to that endpoint by default.
- The **Dart Debug Chrome extension** may be missing, or the browser session isn’t attached to the active VM service.
- Docker networking / port mapping for the VM service can interfere with the tool’s handshake.
- Service worker or caching can sometimes block updates.

**Potential directions (not yet integrated):**
- Use `flutter run -d chrome` (requires a Chrome installation that the Flutter tool can control).
- Map a stable VM service port (e.g., `--host-vmservice-port 5555` → expose with `-p 5555:5555`) and ensure the browser attaches.
- Install/enable the **Dart Debug Extension** in Chrome and confirm it’s connected.
- Hard refresh the app (`Ctrl+Shift+R`), or clear site data to avoid stale service worker caches.

---

## Building an APK (later, same image)
When you want an Android binary:
```bash
make run
# inside the container:
cd /app/ex00
flutter build apk --release
# resulting file:
# build/app/outputs/flutter-apk/app-release.apk
```

This uses the Android SDK components installed by your Dockerfile.

---

## Troubleshooting
- **Port not reachable**: Ensure `-p 8080:8080` is present and the app is running with `--web-hostname 0.0.0.0`.
- **Permissions on host**: If you hit permission errors, ensure your host user can write into `mobilePiscine/`. (Temporary unblock: `chmod -R 777 mobilePiscine`, though a proper UID/GID match is preferred.)
- **`ex00` not found**: Create the app folder with `flutter create ex00` under `/app` (i.e., host `mobilePiscine/`).

---

## Quick command recap
```bash
make build         # Build the image
make web           # Start the web server for Flutter on 0.0.0.0:8080
# Visit http://localhost:8080 on the host
```
