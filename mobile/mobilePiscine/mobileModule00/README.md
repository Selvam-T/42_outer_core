# MuteTimer

MuteTimer is a small Flutter + Android app that lets me mute media volume from a **Quick Settings tile**, wait for a timer to finish, then restore the original volume automatically.

This project is intentionally **Android-first**. The real behavior is implemented on the Android side, while Flutter is kept minimal.

---

## Purpose

The app was created to:

- mute audio quickly during ads in apps like **YouTube** and **Spotify**
- restore the previous media volume automatically after the timer ends
- work from **Quick Settings** without needing a full app UI

---

## Final direction of the project

Older ideas in this project included a Flutter UI with buttons and a timer display.

The newer direction replaced that with this design:

- **no normal Flutter UI flow**
- **minimal `main.dart`**
- **Quick Settings tile** handled by a separate Android `TileService`
- **background restore** handled by Android alarm + receiver logic
- test first with **10 seconds**
- later switch to the real duration such as **3 minutes**

This is the direction this README documents.

---

## High-level architecture

The app has three parts:

1. **App launcher / settings screen**
   - opened when I tap the app icon
   - lets me choose timer duration such as 1 / 2 / 3 minutes
   - shows About information

2. **Quick Settings tile**
   - appears in Android Quick Settings after the app is installed
   - tapping the tile mutes the media stream
   - tile state is managed by Android through the `TileService`

3. **Background restore**
   - Android schedules a restore event
   - when timer ends, original volume is restored
   - tile state returns to inactive

---

## Sequence of work

The order used in this project matters.

### 1. Create the Quick Settings tile icon first

The Quick Settings tile uses its own icon resource.

The tile icon is:

- a **white `VectorDrawable`**
- on a **transparent background**
- stored inside the Android project as a drawable resource

For this project, the tile icon file is:

```text
android/app/src/main/res/drawable/mute_audio.xml
```

This icon is **packaged inside the APK**. It is not copied manually to a folder on the phone.

### 2. Create the normal launcher icon separately

The normal app icon is different from the Quick Settings tile icon.

For this project, the launcher icon source is:

```text
assets/icon/app_icon.png
```

This is used for the app launcher icon, not the Quick Settings tile icon.

### 3. Add Android-side logic

The app flow is Android-side:

- tile service mutes the volume
- receiver restores the volume later
- optional settings screen lets me choose timer length

### 4. Keep Flutter minimal

The Flutter side does not run the mute workflow.

`lib/main.dart` is intentionally minimal:

```dart
import 'package:flutter/widgets.dart';

void main() {
  runApp(const SizedBox.shrink());
}
```

A white blank screen in Chrome is expected from this file and is **not** the Android app behavior.

---

## Project naming

Use these final names consistently.

### Flutter project name

Use:

```text
mute_audio
```

This follows the normal lowercase + underscore naming style.

Do **not** use:

```text
muteAudio
```

### App display name

Use:

```text
MuteTimer
```

### Kotlin package

Use:

```text
com.example.mute_audio
```

### Quick Settings drawable name

Use:

```text
mute_audio.xml
```

### Launcher icon source

Use:

```text
assets/icon/app_icon.png
```

---

## Folder and file layout

Important files in the final project:

```text
mute_audio/
├── lib/
│   └── main.dart
├── assets/
│   └── icon/
│       └── app_icon.png
├── android/
│   └── app/
│       ├── build.gradle.kts
│       └── src/
│           └── main/
│               ├── AndroidManifest.xml
│               ├── kotlin/com/example/mute_audio/
│               │   ├── MainActivity.kt
│               │   ├── MuteTileService.kt
│               │   └── RestoreVolumeReceiver.kt
│               └── res/
│                   ├── drawable/
│                   │   └── mute_audio.xml
│                   ├── layout/
│                   │   └── activity_settings.xml
│                   └── values/
│                       ├── strings.xml
│                       └── styles.xml
├── analysis_options.yaml
├── pubspec.yaml
└── test/
    └── widget_test.dart
```

---

## Why `MuteTileService.kt` exists

`MuteTileService.kt` is created as a **separate Android service class** for the Quick Settings tile.

It is needed because the Quick Settings tile is **not** the same thing as the app’s main activity.

### `MainActivity.kt`

`MainActivity.kt` is the activity opened when I tap the app icon.

In this project, it is used for:

- timer selection
- About information

### `MuteTileService.kt`

`MuteTileService.kt` is the Android `TileService` for the Quick Settings tile.

It is used for:

- handling tile taps
- muting the media stream
- switching tile state between active/inactive
- scheduling the later restore

### They are not the same

- `MainActivity.kt` = app screen
- `MuteTileService.kt` = Quick Settings tile behavior

Do not replace `MainActivity.kt` with tile code.

Create `MuteTileService.kt` as a separate file.

---

## Why `RestoreVolumeReceiver.kt` exists

`RestoreVolumeReceiver.kt` is used to restore the original media volume after the timer ends.

Without it:

- the tile can mute the volume
- but there is no separate component left to restore the old volume later

So `RestoreVolumeReceiver.kt` must remain in the project.

---

## Why `strings.xml` may need to be created manually

Sometimes I may see `styles.xml` but **not** `strings.xml`.

That is not a problem.

If `strings.xml` does not exist, I must create it manually here:

```text
android/app/src/main/res/values/strings.xml
```

### Why it is needed

`strings.xml` stores text resources such as:

- app name
- tile label
- timer labels
- About text

### How it is different from `styles.xml`

- `strings.xml` = text values
- `styles.xml` = appearance / theme / styling rules

So both files can exist together under:

```text
android/app/src/main/res/values/
```

---

## Quick Settings tile behavior

### How the tile is added

After installing the APK, the tile does **not** appear automatically in the visible Quick Settings row.

I must add it manually:

1. swipe down from top
2. expand Quick Settings fully
3. tap edit / pencil / add buttons
4. find the tile
5. drag it into the active tiles area
6. save / done

### What the tile does

When tapped:

1. app reads the current media volume
2. saves the original volume
3. sets media volume to 0
4. sets tile state active
5. schedules a restore event

When timer ends:

1. receiver restores original volume
2. active flag is cleared
3. tile returns to inactive state

### Tile appearance

Android controls the visual treatment of active vs inactive tile state.

On some phones, the background / icon color difference is obvious.  
On other phones, the visual difference may be small.

So if the tile does not visibly invert colors much, that may be device-specific.

### Active tile metadata

The service includes:

```xml
<meta-data
    android:name="android.service.quicksettings.ACTIVE_TILE"
    android:value="true" />
```

This belongs inside the `<service>` block for the tile service.

---

## App uninstall vs tile removal

These are two different actions.

### 1. User removes tile from Quick Settings

This is manual tile removal by the user.

This is what `onTileRemoved()` refers to.

### 2. User uninstalls the app

This is app uninstall.

This is **not** the same as manual tile removal.

Important point:

- `onTileRemoved()` applies to manual tile removal
- do **not** rely on `onTileRemoved()` as an uninstall callback

When the app is uninstalled, the tile provider is gone because the `TileService` is gone.

---

## Timer strategy used during development

For testing, I used a temporary short duration first.

### Testing value

```kotlin
const val TEST_DURATION_MS = 10_000L
```

This means:

- `10_000L` = 10 seconds

### Real duration examples

- 1 minute = `60_000L`
- 2 minutes = `120_000L`
- 3 minutes = `180_000L`

During testing, keeping the duration at 10 seconds makes it much easier to verify:

- mute works
- restore works
- tile resets correctly

After testing, switch to the real duration logic.

---

## Settings screen design

The launcher app screen is intentionally small.

When I tap the app icon, it opens a simple settings page containing:

- timer selection
  - 1 minute
  - 2 minutes
  - 3 minutes
- About section

The selected timer is stored in `SharedPreferences`.

Then `MuteTileService.kt` reads that saved value when the Quick Settings tile is tapped.

---

## About text

Current About text:

```xml
<string name="about_text">MuteTimer was created by Selvam to mute audio during ads in YouTube and Spotify and automatically restore the previous volume after the selected timer ends. To activate the timer, add the Mute Timer tile to Quick Settings, then tap the tile whenever you want to mute audio temporarily.</string>
```

---

## Files to create or modify

These are the main files needed.

### Create

- `assets/icon/app_icon.png`
- `android/app/src/main/res/drawable/mute_audio.xml`
- `android/app/src/main/kotlin/com/example/mute_audio/MuteTileService.kt`
- `android/app/src/main/kotlin/com/example/mute_audio/RestoreVolumeReceiver.kt`
- `android/app/src/main/res/layout/activity_settings.xml`
- `android/app/src/main/res/values/strings.xml`

### Modify

- `lib/main.dart`
- `pubspec.yaml`
- `android/app/src/main/kotlin/com/example/mute_audio/MainActivity.kt`
- `android/app/src/main/AndroidManifest.xml`
- `android/app/build.gradle.kts`
- `analysis_options.yaml`

### Replace or simplify

- `test/widget_test.dart`

---

## Final file contents summary

### `lib/main.dart`

Keep it minimal:

```dart
import 'package:flutter/widgets.dart';

void main() {
  runApp(const SizedBox.shrink());
}
```

### `MainActivity.kt`

Use it for:

- selecting timer duration
- showing About text

### `MuteTileService.kt`

Use it for:

- tile tap handling
- muting audio
- scheduling restore
- updating tile state

### `RestoreVolumeReceiver.kt`

Use it for:

- restoring original volume
- resetting active state

### `activity_settings.xml`

Use it for:

- radio buttons for 1 / 2 / 3 minutes
- About section

### `strings.xml`

Use it for:

- `app_name`
- `tile_label`
- timer labels
- About text

### `AndroidManifest.xml`

Use it to declare:

- launcher activity
- tile service
- restore receiver
- exact alarm permission if required

### `build.gradle.kts`

Use it for the app module configuration.

Do not confuse:

- top-level `build.gradle.kts`
- `android/app/build.gradle.kts`

App-specific dependencies belong in:

```text
android/app/build.gradle.kts
```

---

## Building the app from scratch

## 1. Create the Flutter project

```bash
flutter create mute_audio
```

## 2. Enter the project

```bash
cd mute_audio
```

## 3. Add or replace project files

Create or update the Android-side files listed above.

## 4. Add assets

Place launcher icon here:

```text
assets/icon/app_icon.png
```

Place Quick Settings vector drawable here:

```text
android/app/src/main/res/drawable/mute_audio.xml
```

## 5. Update `pubspec.yaml`

Make sure the launcher asset is included and keep only one `flutter:` block.

Example structure:

```yaml
name: mute_audio
description: "A Flutter project."
publish_to: 'none'

version: 1.0.0+1

environment:
  sdk: ^3.9.0

dependencies:
  flutter:
    sdk: flutter

dev_dependencies:
  flutter_test:
    sdk: flutter
  flutter_lints: ^5.0.0
  flutter_launcher_icons: ^0.14.4

flutter_launcher_icons:
  android: true
  ios: false
  image_path: "assets/icon/app_icon.png"

flutter:
  uses-material-design: true
  assets:
    - assets/icon/app_icon.png
```

Important:
- do **not** create duplicate `flutter:` keys
- `flutter_launcher_icons:` is a separate top-level block

## 6. Get packages

```bash
flutter pub get
```

## 7. Generate launcher icon if using launcher icon package

If I use `flutter_launcher_icons`, run:

```bash
dart run flutter_launcher_icons
```

## 8. Analyze project

```bash
flutter analyze
```

## 9. Build APK

```bash
flutter build apk --release
```

## 10. Install APK on phone

Install the built APK on the Android phone.

---

## Where the built APK usually appears

Usually under:

```text
build/app/outputs/flutter-apk/app-release.apk
```

---

## Installing on the phone

If installing from outside Play Store, Android may require allowing installs from that source.

If I move the APK using Google Drive:

1. upload APK to Google Drive
2. download APK on the phone
3. allow install from that source if Android asks
4. install the APK

---

## After installation

1. open the app icon once if needed
2. choose timer duration
3. return to home
4. open Quick Settings edit mode
5. add the Mute Timer tile
6. tap the tile to test mute behavior

---

## Testing workflow

Recommended test order:

### 1. Use 10 seconds first

Test that:

- tile mutes volume
- volume restores correctly
- tile state resets

### 2. Then switch to real timer values

After confirmation, use:

- 1 min
- 2 min
- 3 min

---

## Web and Chrome behavior

If I run:

```bash
flutter run
```

and Flutter opens **Chrome** with a blank white page, that is expected for this project.

Why:
- Flutter is launching the **web target**
- `main.dart` intentionally renders nothing
- this does **not** test the Android app flow

So:
- blank white page in Chrome = expected
- not evidence that Android app is broken

---

## ADB and device testing notes

If:

```bash
adb logcat
```

says:

```text
waiting for device
```

that means `adb` does not currently see an Android device or emulator.

So if `flutter run` opens in Chrome and `adb` says waiting for device:

- I am testing the web target
- I am **not** testing the Android app crash path

---

## Common issues and fixes

### 1. Duplicate `flutter:` key in `pubspec.yaml`

Symptom:
- YAML parse error
- duplicate mapping key

Fix:
- keep only one top-level `flutter:` block

### 2. Drawable not found

Symptom:
- manifest references drawable not found

Fix:
- make sure file really exists at:

```text
android/app/src/main/res/drawable/mute_audio.xml
```

- manifest must match the exact file name:

```xml
android:icon="@drawable/mute_audio"
```

### 3. Default widget test still refers to old app class

Symptom:
- `widget_test.dart` expects a widget class that no longer exists

Fix:
- replace with a placeholder test
- or remove the old generated test content

### 4. `flutter_lints` include warning

Fix options:
- add `flutter_lints` under `dev_dependencies`
- or simplify `analysis_options.yaml`

### 5. Tile mutes but does not restore

Possible causes:
- exact alarm permission / special access issue
- restore receiver not declared
- alarm not firing on that device
- testing with Android behavior differences

### 6. Tile color does not visibly change much

Possible cause:
- OEM / device-specific Quick Settings appearance

---

## Recommended final polish steps

After the core app works:

1. switch from 10-second test duration to selected timer values
2. confirm the settings screen saves selected duration
3. rebuild release APK
4. test on the real device again
5. keep names consistent:
   - project = `mute_audio`
   - app name = `MuteTimer`
   - tile icon = `mute_audio.xml`

---

## Final notes

This app is built around Android Quick Settings, not around a normal Flutter UI.

That is why:

- Flutter `main.dart` is minimal
- the mute workflow is handled in native Android code
- tile behavior lives in `MuteTileService.kt`
- timed restore lives in `RestoreVolumeReceiver.kt`
- launcher screen is only for settings/about

That separation is intentional and is the correct mental model for this app.
