# Bare Metal Installation
- ensure you have `docker` and `go` installed
- clone the repository
```bash
git clone https://github.com/HarshVaragiya/LambdaFn
```
- install the binary using the makefile
```bash
make install
```
- Makefile builds lambdaRuntime which is used to communicate with docker containers, then builds the runtime container images, and then finally builds and installs the LambdaFn binary
- run binary from ~/go/bin/LambdaFn
```bash
cd ~/go/bin/
./LambdaFn
```