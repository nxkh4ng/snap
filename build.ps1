$targets = @(
    @{ GOOS="windows"; GOARCH="amd64"; OUT="snap_Windows_x86_64.exe" },
    @{ GOOS="darwin";  GOARCH="amd64"; OUT="snap_Darwin_x86_64" },
    @{ GOOS="darwin";  GOARCH="arm64"; OUT="snap_Darwin_arm64" },
    @{ GOOS="linux";   GOARCH="amd64"; OUT="snap_Linux_x86_64" }
)

New-Item -ItemType Directory -Force -Path ./bin | Out-Null

foreach ($t in $targets) {
    $env:GOOS = $t.GOOS
    $env:GOARCH = $t.GOARCH
    go build -o "./bin/$($t.OUT)" .
    Write-Host "Built $($t.OUT)"
}

Remove-Item Env:GOOS
Remove-Item Env:GOARCH

Write-Host "Done. Binaries in ./bin/"
