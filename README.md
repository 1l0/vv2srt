# VOICEVOX to SRT

## Download

[Latest release](https://github.com/1l0/voicevox2srt/releases/latest)

## Usage

1. **VOICEVOX or AivisSpeech**: Export or preview all of audio (this updates the necesarry parameters in your project).
2. **VOICEVOX or AivisSpeech**: Save the project.
3. Run the following command with the saved project file path.
    - By default, `<project file path>.srt` should be generated.

### VOICEVOX

```sh
voicevox2srt <project name>.vvproj
```

### AivisSpeech

```sh
voicevox2srt <project name>.aisp
```

### Options

Specify the output file path:

```sh
voicevox2srt -o example.srt <project name>.vvproj
```
