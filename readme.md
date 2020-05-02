# Extractor

Opinionated zip/rar extractor for torrent clients

![Latest Release](https://github.com/napei/extractor/workflows/Upload%20Release%20Asset/badge.svg)

---

## Download

Download latest release here: [Latest Release](https://github.com/napei/extractor/releases/latest)

## Commandline usage

Application can be run from the command line. Typically an input directory is specified using the `-input` flag, and then the directory is searched recursively and all `.rar`,`.part01.rar`, and `.zip` files are extracted and their contents are left in the same directory as the archive.

This can be applicable to many situations, however this executable was designed specifically to work well with torrents, as 99% of torrents that have archives will either have a master `.rar` file or many `.part0x.rar` files. This application will only look at `.part01.rar` files in the latter case, so as to not add extra cost to extraction.

The output directory can be configured through the `-output` flag, and all items will be extracted to that directory instead of the directory of the archive.

The application can also perform a 'dry run' through the use of the `-dryrun` flag, allowing for archives to be found and counted/listed but not extracted.

By default, the application will hide names of the files and only show a count of the archives found and a status with a number, to prevent visual spam. These messages can be enabled through the `-verbose` flag.

Also, if an output file already exists, the application will overwrite the file by default. This can be overridden with the `-no-overwrite` flag, which will disalow overwriting existing files.

Basic terminal usage is shown below, where the filename relates to the latest version downloaded.

`> extractor.v0.x.exe -input C:\Torrents\Seeding\`

## Usage in torrent clients

Most modern torrent clients allow for the running of a program on torrent completion. This feature allows arguments relevant to that torrent to be passed to the program also. Examples are shown below.

### qBittorrent

qBittorrent is known to have issues with executing a program when torrents are in automatic management mode, as the execution is on torrent completion and not when the torrent finishes moving. Depending on your system, the timeout value shown below may need to be adjusted to suit. If the downloading and seeding directories are on the same drive, 10 seconds is plenty.

![qBittorrent Usage Example](https://raw.githubusercontent.com/napei/extractor/master/images/qBittorrent%20Usage.png)
