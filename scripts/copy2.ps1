# Copy WVD -> Server

$SRC_BIN="\\tsclient\c\temp1\checkmirror\checkmirror.exe"
$DST_FOLDER="c:\checkmirror"

mkdir -Force "$DST_FOLDER"
copy "$SRC_BIN" "$DST_FOLDER"
