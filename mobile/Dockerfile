# 1. THE BASE (Linux + Basics)
FROM ubuntu:22.04

# Prevent interactive prompts (like "What is your time zone?")
ENV DEBIAN_FRONTEND=noninteractive

# Install basics: curl, git, unzip, and Java (JDK 17)
# 'ca-certificates' is crucial for SSL as discussed before
RUN apt-get update && apt-get install -y \
    curl \
    git \
    unzip \
    xz-utils \
    zip \
    libglu1-mesa \
    openjdk-17-jdk \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# 2. THE ENVIRONMENT VARIABLES
# We set up where Android and Flutter will live
ENV ANDROID_HOME=/usr/local/android-sdk
ENV FLUTTER_HOME=/usr/local/flutter
ENV PATH="$ANDROID_HOME/cmdline-tools/latest/bin:$ANDROID_HOME/platform-tools:$FLUTTER_HOME/bin:$PATH"

# 3. ANDROID SDK SETUP
# We download the "Command Line Tools" only.
# (Note: This URL changes occasionally, check developer.android.com/studio)
WORKDIR /tmp
RUN curl -O https://dl.google.com/android/repository/commandlinetools-linux-11076708_latest.zip

# Create the folder structure. 
# IMPORTANT: sdkmanager usually requires the folder structure: cmdline-tools/latest/bin
RUN mkdir -p $ANDROID_HOME/cmdline-tools \
    && unzip commandlinetools-linux-*_latest.zip -d $ANDROID_HOME/cmdline-tools \
    && mv $ANDROID_HOME/cmdline-tools/cmdline-tools $ANDROID_HOME/cmdline-tools/latest \
    && rm commandlinetools-linux-*_latest.zip

# 4. INSTALL SPECIFIC ANDROID COMPONENTS
# We use 'yes' to automatically accept licenses.
# We install exactly what you listed: build-tools 34, platform 34, etc.
RUN yes | sdkmanager --licenses \
    && sdkmanager "platform-tools" "platforms;android-34" "build-tools;34.0.0"

# 5. PREPARE NON-ROOT USER 
# We create a user named 'developer' with UID 1000 (Standard linux user_
RUN useradd -m -s /bin/bash -u 1000 developer

# We give this user ownership of the SDK folders so they can update them
RUN mkdir -p $FLUTTER_HOME
RUN chown -R developer:developer $ANDROID_HOME $FLUTTER_HOME

# 6. FLUTTER SETUP
# switch to a new user BEFORE downloading Flutter
USER developer
WORKDIR /home/developer

# We clone the stable version of Flutter
RUN git clone https://github.com/flutter/flutter.git $FLUTTER_HOME -b stable

# Tell git to trust this directory (fixes "dubious ownership" errors)
RUN git config --global --add safe.directory $FLUTTER_HOME

# Disable analytics and precache (downloads basic binaries)
RUN flutter config --no-analytics \
    && flutter config --enable-android \
    && flutter precache

# 6. WORKSPACE
# Create a folder for your app code
WORKDIR /app
