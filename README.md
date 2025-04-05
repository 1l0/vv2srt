# VOICEVOX to SRT

## Download

[Latest release](https://github.com/1l0/voicevox2srt/releases/latest)

## Usage

1. **VOICEVOX**: Export or preview all of audio (this updates the necesarry parameters in your project).
2. **VOICEVOX**: Save the project.
3. Run the following command with the saved project file path.
    - By default, `subtitles.srt` should be generated in the current directory.

### VOICEVOX

```sh
voicevox2srt project_name.vvproj
```

### AivisSpeech

```sh
voicevox2srt project_name.aisp
```

### Options

Specify the output file path:

```sh
voicevox2srt -o example.srt some_project.vvproj
```
