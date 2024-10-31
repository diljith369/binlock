# BinLock 🔒
BinLock is a cross-platform security tool designed to add password protection to executable files. It consists of two main components: `binguard` for protecting binaries and `launcher` for executing protected binaries.

## Features ✨

- Password protect any executable file
- Cross-platform compatibility (Windows, Linux, macOS)
- Secure password verification
- Simple and straightforward usage
- Minimal performance impact

## Components 🛠️

### BinGuard
The protection utility that encrypts and secures your binary files.

### Launcher
The execution utility that runs protected binaries after password verification. If the protected binary has command line arguments you can provide it at the end of the launcher after password.

## Usage  💻

### Protecting a Binary

To protect an executable file using BinGuard:

binguard -i <input_binary> -o <output_binary> -pass <password>

Example:

binguard.exe -i .\Presentation.exe -o protected.exe -pass p@ssword

![binlock](https://github.com/diljith369/binlock/blob/main/binguard.png)

### Running a Protected Binary

To run a protected binary using Launcher:

launcher <protected-binary> <password> [binary args...]

Example:

launcher.exe prbin.exe p@ssword cmdargs1 cmdargs2

## Command Line Arguments

### BinGuard
- -i: Input binary path
- -o: Output binary path
- -pass: Password to protect the binary

### Launcher
- -ipf: Protected binary path
- -pass: Password to unlock and run the binary

## Security Considerations 🛡️

- Choose strong passwords
- Keep passwords secure and confidential
- BinLock is not a replacement for code signing or other security measures
- Protected binaries should be handled with the same security precautions as unprotected ones

## Disclaimer ⚠️
This tool is provided as-is without any warranties and mainly for educational purposes. Users are responsible for their use of this software and should ensure compliance with local laws and regulations.


