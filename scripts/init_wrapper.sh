#!/bin/bash

docker exec -it gymdl /app/wrapper/wrapper -L '$1':'$2' -F -H 0.0.0.0

#docker cp gymdl:/app/wrapper ~/data/gymdl/wrapper