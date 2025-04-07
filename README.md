# vv2srt

This command generates an SRT (SubRip) from a VOICEVOX or AivisSpeech project.

## Download

[Latest release](https://github.com/1l0/vv2srt/releases/latest)

## Usage

1. **VOICEVOX or AivisSpeech**: Export or preview all of audio in your project.
    - This updates the necesarry parameters in the project.
2. **VOICEVOX or AivisSpeech**: Save the project.
3. Run the following command with the saved project file path.
    - By default, `<project file path>.srt` should be generated.

### VOICEVOX

```sh
vv2srt <project name>.vvproj
```

### AivisSpeech

```sh
vv2srt <project name>.aisp
```

### Options

Specify the output file path:

```sh
vv2srt -o example.srt <project name>.vvproj
```
