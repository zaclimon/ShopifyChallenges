# Memory Game 🃏

This game is a sample application replicating what the game of [Concentration](https://en.wikipedia.org/wiki/Concentration_(card_game)) for Android devices. This has been made as part of [Shopify](https://www.shopify.com/)'s Coding Challenge for Winter 2020 internships. The goal here is to represent a regular game without too much fluff mostly due to time constraints. 😅

**Note**: This application requires internet access in order to download the game's assets on it's first run.

## Libraries Used

- [OkHttp](https://square.github.io/okhttp/) for streamlined HTTP based requests
- [Glide](https://bumptech.github.io/glide/) for simpler image processing
- [Moshi](https://github.com/square/moshi) for easier JSON processing
- [Material Components](https://material.io/develop/android/) for some styling

## SDK Components

- Target API: [29 (Android 10)](https://developer.android.com/studio/releases/platforms#10)
- Minimum API: [24 (Android 7.0 Nougat)](https://developer.android.com/studio/releases/platforms#7.0)
- [Android SDK Tools 26.1.1](https://developer.android.com/studio/releases/sdk-tools#notes)
- [Android Platform Tools 29.0.4](https://developer.android.com/studio/releases/platform-tools)
- [Android Emulator 29.2.0](https://androidstudio.googleblog.com/2019/09/emulator-2920-stable.html)

## Devices Tested

- Android Emulator using a Pixel 2 frame
- [Google Pixel](https://en.wikipedia.org/wiki/Pixel_(smartphone))

Both devices were running on [Android 10](https://en.wikipedia.org/wiki/Android_10) when testing the application.

## Permissions

- [`android.permission.INTERNET`](https://developer.android.com/reference/android/Manifest.permission.html#INTERNET) for downloading assets when the application is started for the first time.
- [`android.permission.ACCESS_NETWORK_STATE`](https://developer.android.com/reference/android/Manifest.permission.html#ACCESS_NETWORK_STATE) for verifying whether there is connectivity for downloading the assets.

## Getting Started

1. Clone this repository

    `git clone https://github.com/zaclimon/MemoryGame`
2. Open in Android Studio
3. Download missing SDK components (if any)
4. Compile
5. Enjoy! 🎉