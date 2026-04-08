# Flutter in Docker:  
## Full Setup vs Web-Only Setup

This README explains the **two Docker/Makefile setups** in this folder, what each one is for, how they differ, how to initialize a Flutter exercise correctly, and how to run them.

Your project currently has **two pairs** of files:

1. **Full setup**
   - `Dockerfile`
   - `Makefile`

2. **Web-only setup**
   - `Dockerfile.web-only`
   - `Makefile.web-only`

The full setup is heavier and meant for **Flutter web now, with Android build support later**.
The web-only setup is lighter and meant for **immediate browser-based development and testing only**.

---

## Separate host folders

You decided to keep **two separate environments** so the two setups do not write into the same mounted host volume.

Current host folders:

- **Full setup** uses:
  ```text
  $(PWD)/mobilePiscine
  ```
- **Web-only setup** uses:
  ```text
  $(PWD)/mobilePiscineWeb
  ```

This prevents one setup from overwriting the other setup's files.

---

## Expected Flutter app paths

The Makefiles are meant to work with a **module root** plus a selectable **exercise folder**.

Default variable pattern:

```make
WORKSPACE_DIR := mobileModule00
EXERCISE := ex00
APP_DIR := $(WORKSPACE_DIR)/$(EXERCISE)
```

That means the Flutter project is expected inside the selected exercise folder, and that folder should contain `pubspec.yaml`.

For the **full setup**, the default Flutter project path is:

```text
mobilePiscine/mobileModule00/ex00
```

For the **web-only setup**, the default Flutter project path is:

```text
mobilePiscineWeb/mobileModule00/ex00
```

If you override the exercise at runtime, the path changes accordingly. Example:

```bash
make EXERCISE=ex01 web
make -f Makefile.web-only EXERCISE=ex01 web
```

Those commands make the app path resolve to:

```text
mobilePiscine/mobileModule00/ex01
```

or:

```text
mobilePiscineWeb/mobileModule00/ex01
```

Inside the container, the host folders are mounted and used as the project workspace, while `APP_DIR` points to the specific Flutter app you want to run.

---

## Initialize the project

Do **not** manually create an empty `ex00` folder and then try to run Flutter inside it.

An empty folder is **not** a Flutter project. Flutter expects the selected app folder to already contain files such as:

```text
pubspec.yaml
lib/
web/
```

The correct approach is:

1. create the Docker image
2. open a shell inside the container
3. go to the module folder
4. run `flutter create ex00`

That command generates the actual Flutter project inside `ex00`.

### Full setup initialization

1, Before creating the first Flutter exercise project, start the container first:

```bash
make build
make run
```
2, Then, in another terminal, run:

```bash
make chmod
```
#### Why make chmod is needed  

In the full setup, the project folder from the host machine is mounted into the container.
The host user and the container user may not have the same UID/GID, so even if the folder exists, the container user may not have permission to write into it.

The make chmod target fixes that by running a permission update inside the running container:

it finds the running container created from the current image
it executes chmod -R a+rwx /app/mobileModule00 as root inside the container
because /app/mobileModule00 is a mounted host folder, the permission change also affects that host folder

This makes the mounted project directory writable.  

3, After that, return to the container terminal and initialize the Flutter project:

```bash
cd /app
mkdir -p mobileModule00
cd mobileModule00
flutter create ex00
```

That creates the Flutter project at:

```text
mobilePiscine/mobileModule00/ex00
```

### Web-only setup initialization

If the file is literally named `Makefile.web-only`, run it like this:

```bash
make -f Makefile.web-only build
make -f Makefile.web-only shell
```

You can also use:

```bash
make -f Makefile.web-only run
```

`make shell` and `make run` both open an interactive container shell.
For initialization, `make shell` is usually more convenient because it starts in `/workspace`.

`make chmod` is **not** part of normal web-only initialization.
Use it only if you hit a mounted-folder permission problem.

Then inside the container:

```bash
cd /workspace
mkdir -p mobileModule00
cd mobileModule00
flutter create ex00
```

That creates the Flutter project at:

```text
mobilePiscineWeb/mobileModule00/ex00
```

### Why not manually create `ex00`

If you manually create:

```text
mobileModule00/ex00
```

as an empty folder, then `make web` or `make web-build` will fail because Flutter will not find a project root in that folder.

---

## 1) Full setup: `Dockerfile` + `Makefile`

**Note:** the full-setup `Makefile` directives may still need some tweaking depending on how you want to organize your Android and web workflow.

### Purpose
Use this set when you want a **broader Flutter environment**.

This setup includes:
- Flutter SDK
- Dart SDK
- Java JDK
- Android command-line tools
- Android SDK components
- NDK

This is the setup to keep when you want to:
- run Flutter in Docker now
- test on web
- later build an Android APK from the same environment
- later explore USB passthrough, device testing, or other Android-related tooling

### Trade-off
- **Pros:** more complete, Android-ready
- **Cons:** much larger image, slower build, more storage used

### Build the full image
```bash
make build
```

### Open a shell inside the full container
```bash
make run
```

### Start the Flutter app on the web
```bash
make web
```

Then open:

```text
http://localhost:8080
```

### Build Android APK with the full setup
```bash
make apk
```

Or use the wrapper target:

```bash
make launch_apk
```

### Launch web with wrapper target
```bash
make launch_web
```

### Other full-setup commands
```bash
make clean
make prune
make chmod
```

### Notes about the full setup
- This is the **default** setup when you run plain `make ...` because the default file name is `Makefile`.
- This setup mounts:
  ```text
  $(PWD)/mobilePiscine -> /app
  ```
- `make web` in this setup may expose extra ports for VM service experiments, depending on your Makefile.
- This setup is the better choice when you eventually want Android build support without redesigning the container.

---

## 2) Web-only setup: `Dockerfile.web-only` + `Makefile.web-only`

### Purpose
Use this set when you want the **smallest practical Flutter container for browser testing**.

This setup includes only what you need for:
- Flutter SDK
- Dart SDK
- Flutter web support
- running the app in a browser through `web-server`
- building a web output with `flutter build web`

It does **not** include Android SDK / JDK / NDK.

### Trade-off
- **Pros:** smaller, cleaner, faster to build, less storage use
- **Cons:** no Android APK build from this image

### Why this is your immediate-purpose setup
This is the right set when your goal is:
- learn Flutter and Dart
- test the app in a browser
- avoid Android emulator or device complexity for now
- reduce Docker image size

### Important difference from the full setup
The web-only setup mounts a **different host folder**:

```text
$(PWD)/mobilePiscineWeb
```

That means it is isolated from the full setup.

---

## How to run the web-only Makefile

If the file is literally named `Makefile.web-only`, tell `make` which file to use.

Example:

```bash
make -f Makefile.web-only build
```

For convenience, you may rename `Makefile.web-only` to `Makefile` before using it day to day.
In the command examples below, the web-only file is assumed to be the active file named `Makefile`.

### Build the web-only image
```bash
make build
```

### Open a shell inside the web-only container
```bash
make shell
```

### Alternate interactive shell target
```bash
make run
```

`make run` also opens a container shell.
For normal development, you still use `make web` on the host.

### Check the Flutter web environment
```bash
make doctor
```

### Run the Flutter app in web-server mode
```bash
make web
```

Then open:

```text
http://localhost:8080
```

### Run a different exercise with `EXERCISE=...`
If your Makefile uses:

```make
WORKSPACE_DIR := mobileModule00
EXERCISE := ex00
APP_DIR := $(WORKSPACE_DIR)/$(EXERCISE)
```

then `EXERCISE := ex00` is only the **default**.

You can override it from the command line for a single run:

```bash
make EXERCISE=ex01 web
```

Why this works:
- `EXERCISE=ex01` overrides the default `ex00` for that command only
- `APP_DIR` becomes `mobileModule00/ex01`
- `make` then runs Flutter commands inside `mobilePiscineWeb/mobileModule00/ex01`
- this is useful because `mobileModule00` can contain multiple exercises such as `ex00`, `ex01`, `ex02`, and `ex03`

Important:
- `make EXERCISE=ex01 ...` does **not** create `ex01` automatically
- the selected exercise folder should already be the actual Flutter app folder
- that means it should contain `pubspec.yaml`

### Build deployable web output
```bash
make web-build
```

The generated output will be in the selected exercise folder. With the default `EXERCISE := ex00`, that is:

```text
mobilePiscineWeb/mobileModule00/ex00/build/web
```

If you run, for example:

```bash
make EXERCISE=ex01 web-build
```

then the output will be:

```text
mobilePiscineWeb/mobileModule00/ex01/build/web
```

### Permissions troubleshooting
```bash
make run
```

In another terminal:

```bash
make chmod
```

Use `make chmod` only when the mounted project folder has permission issues.
It is not required as a normal initialization step.

### Remove the web-only image
```bash
make clean
```

### Remove image and prune Docker cache
```bash
make prune
```

---

## Which set should you use?

### Use the full setup (`Dockerfile` + `Makefile`) when:
- you want Android build support
- you want one container that can later build APKs
- you want to explore device testing later
- you are okay with a heavier Docker image

### Use the web-only setup (`Dockerfile.web-only` + `Makefile.web-only`) when:
- you only want browser testing right now
- you want less storage use
- you want a simpler container
- Android or USB passthrough work can wait until later

For your current goal, **use the web-only setup**.

---

## Recommended workflow right now

### First time
```bash
make build
make shell
```

Then inside the container:

```bash
cd /workspace
mkdir -p mobileModule00
cd mobileModule00
flutter create ex00
```

### Normal web development
```bash
make web
```

Then visit:

```text
http://localhost:8080
```

Only use `make chmod` if you later run into a permission problem with files under `mobilePiscineWeb`.

### When you want a build artifact for the browser
```bash
make web-build
```

---

## Summary

- `Dockerfile` + `Makefile` = heavier environment in `mobilePiscine`
- `Dockerfile.web-only` + `Makefile.web-only` = lighter environment in `mobilePiscineWeb`
- keeping separate host folders avoids overwrite conflicts
- do **not** manually create an empty exercise folder and treat it as a Flutter project
- initialize the exercise with `flutter create ex00`
- use `EXERCISE=...` to switch between `ex00`, `ex01`, `ex02`, and `ex03`
- `make run` is an alternate interactive shell target in the web-only setup
- `make chmod` in the web-only setup is for permission troubleshooting only, not normal initialization
- for now, your **web-only setup is the better day-to-day workflow**
