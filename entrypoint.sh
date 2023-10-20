#!/bin/sh

while ! nc -z server 8901; do
  echo "Waiting for server to be ready..."
  sleep 1
done

# Now, start the client application
/app/client upload -f 1.txt -f 2.txt
/app/client download -i 0a23afe0bd19853e6c11f581cc6c59afe1fc6b208b471011cf2338a91e4094ba
/app/client download -i 1c46e2f0f5767113dff10781f257ac87a8163c09a201d6bbc466ab6e302ff2fe

# Check if files are downloaded successfully
for file in 0a23afe0bd19853e6c11f581cc6c59afe1fc6b208b471011cf2338a91e4094ba 1c46e2f0f5767113dff10781f257ac87a8163c09a201d6bbc466ab6e302ff2fe; do
  if [[ -f "/app/$file" ]]; then
    echo "OK: $file has been downloaded successfully!"
  else
    echo "Err: $file was not downloaded."
  fi
done

# Keep the container running (you can remove this if you want the container to exit after the checks)
tail -f /dev/null