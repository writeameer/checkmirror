# Copy Laptop->WVD

$SRC_BIN="\\tsclient\c\users\ameer\progs\checkmirror\checkmirror.exe"
$WVD_BIN="c:\temp1\checkmirror\checkmirror.exe"


if (Test-Path($WVD_BIN)) {
  Remove-Item -force "$WVD_BIN" | Out-Null
} 

Copy-Item "$SRC_BIN" "$WVD_BIN"