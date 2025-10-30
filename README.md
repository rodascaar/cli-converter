# audioconv

CLI profesional de conversión de audio multi-formato en Go

[![Go Report Card](https://goreportcard.com/badge/github.com/example/audioconv)](https://goreportcard.com/report/github.com/example/audioconv)
[![GoDoc](https://godoc.org/github.com/example/audioconv?status.svg)](https://godoc.org/github.com/example/audioconv)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Herramienta de línea de comandos escrita en Go para convertir archivos de audio entre múltiples formatos con control fino sobre parámetros de calidad, optimizada para uso en terminal con progress bars, validaciones y salida elegante.

## Características

- **Conversión simple**: `audioconv convert input.flac output.mp3`
- **Conversión en lote**: `audioconv batch *.flac --to mp3`
- **Auto-detección de formato de entrada**
- **Preservación de metadatos (ID3 tags)** cuando sea posible
- **Control fino de calidad**: bitrate, sample rate, canales, presets
- **Progress bars elegantes** con ETA y velocidad
- **Validaciones robustas** pre y post conversión
- **Sistema de presets**: low, medium, high, ultra, lossless
- **Normalización de volumen** a -16 LUFS
- **Trimming**: recortar inicio/fin de archivos
- **Configuración persistente** en `~/.audioconv.yaml`

## Instalación

### Desde código fuente

```bash
git clone https://github.com/example/audioconv.git
cd audioconv
make build
sudo make install
```

### Usando Go

```bash
go install github.com/example/audioconv@latest
```

### Requisitos

- **Go 1.21+**
- **FFmpeg** instalado y en PATH

```bash
# Ubuntu/Debian
sudo apt install ffmpeg

# macOS
brew install ffmpeg

# Windows (usando Chocolatey)
choco install ffmpeg
```

## Uso Rápido

### Conversión básica

```bash
# Convertir un archivo FLAC a MP3 de alta calidad
audioconv convert album.flac album.mp3 -p high

# Convertir con parámetros específicos
audioconv convert song.wav song.mp3 -b 320k -q 0
```

### Conversión en lote

```bash
# Convertir todos los FLAC a MP3
audioconv batch '*.flac' --to mp3 -b 320k

# Conversión paralela con 8 procesos
audioconv batch 'music/**/*.wav' --to opus -p high -j 8
```

### Información de archivos

```bash
# Ver información detallada
audioconv info song.mp3

# Salida JSON
audioconv info music.flac --json

# Solo un campo específico
audioconv info podcast.m4a -f duration
```

## Comandos

### `convert`

Convierte un solo archivo de audio.

```bash
audioconv convert <input> <output> [flags]
```

**Flags:**
- `-b, --bitrate string`: bitrate de salida (ej: 192k, 320k)
- `-q, --quality int`: calidad VBR 0-9 para mp3/ogg (0=mejor)
- `-s, --sample-rate int`: frecuencia de muestreo (ej: 44100, 48000)
- `-c, --channels string`: canales: mono|stereo|auto (default "auto")
- `-p, --preset string`: preset de calidad: low|medium|high|ultra|lossless
- `--normalize`: normalizar volumen a -16 LUFS
- `--trim-start string`: recortar inicio (ej: 10s, 1m30s)
- `--trim-end string`: recortar desde el final
- `-y, --overwrite`: sobrescribir sin preguntar
- `--dry-run`: mostrar comando ffmpeg sin ejecutar
- `--quiet`: suprimir output excepto errores
- `-v, --verbose`: mostrar salida completa de ffmpeg

### `batch`

Convierte múltiples archivos que coincidan con un patrón.

```bash
audioconv batch <pattern> --to <format> [flags]
```

**Flags adicionales:**
- `--to string`: formato de destino (requerido)
- `-o, --output-dir string`: directorio de salida
- `-j, --parallel int`: número de conversiones simultáneas (default 4)
- `-r, --recursive`: buscar archivos recursivamente
- `--keep-structure`: mantener estructura de carpetas

### `info`

Muestra información detallada del archivo de audio.

```bash
audioconv info <file> [flags]
```

**Flags:**
- `--json`: salida en formato JSON
- `-f, --format string`: mostrar solo un campo (duration, bitrate, codec)

## Presets de Calidad

| Preset | MP3 | Opus | AAC | Uso recomendado |
|--------|-----|------|-----|-----------------|
| `low` | 128k | 64k | 128k | podcasts, voz |
| `medium` | 192k | 96k | 192k | música casual |
| `high` | 320k | 128k | 256k | música high quality |
| `ultra` | VBR 0 | 160k | 320k | archival |
| `lossless` | - | - | - | máxima calidad |

## Configuración

Crea un archivo `~/.audioconv.yaml` para valores por defecto:

```yaml
defaults:
  output_format: mp3
  bitrate: 192k
  quality: 2
  sample_rate: 44100
  channels: auto
  normalize: false
  overwrite: false
  output_dir: ./converted
  parallel_jobs: 4

presets:
  podcast:
    format: mp3
    bitrate: 64k
    channels: mono
    normalize: true
  hifi:
    format: flac
    compression_level: 8
    sample_rate: 96000

ffmpeg:
  path: ffmpeg
  extra_args: []
```

## Formatos Soportados

### Entrada
- MP3, AAC, FLAC, WAV, OGG Vorbis, Opus, WMA, ALAC

### Salida
- MP3, AAC, FLAC, WAV, OGG Vorbis, Opus

## Ejemplos de Uso

### Música de alta calidad
```bash
audioconv convert album.flac album.mp3 -p high
```

### Podcast optimizado
```bash
audioconv convert episode.wav episode.mp3 -c mono -b 64k --normalize
```

### Conversión masiva
```bash
audioconv batch '*.flac' --to opus -b 96k -o ./opus_output
```

### Archival
```bash
audioconv convert master.wav master.flac -p lossless
```

## Desarrollo

```bash
# Clonar repositorio
git clone https://github.com/example/audioconv.git
cd audioconv

# Instalar dependencias
go mod tidy

# Ejecutar tests
make test

# Build
make build

# Desarrollo con recarga automática
make dev
```

## Contribuir

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## Licencia

Este proyecto está bajo la Licencia MIT - ver el archivo [LICENSE](LICENSE) para más detalles.

## Soporte

- Reporta bugs en [GitHub Issues](https://github.com/example/audioconv/issues)
- Preguntas en [GitHub Discussions](https://github.com/example/audioconv/discussions)

---

Hecho con ❤️ en Go