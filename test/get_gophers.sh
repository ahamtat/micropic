#!/bin/bash

curl http://localhost:8080/fill/50/50/imagesource/_gopher_original_1024x504.jpg     --output gopher_50x50.jpg
curl http://localhost:8080/fill/200/700/imagesource/_gopher_original_1024x504.jpg   --output gopher_200x700.jpg
curl http://localhost:8080/fill/256/126/imagesource/_gopher_original_1024x504.jpg   --output gopher_256x126.jpg
curl http://localhost:8080/fill/333/666/imagesource/_gopher_original_1024x504.jpg   --output gopher_333x666.jpg
curl http://localhost:8080/fill/500/500/imagesource/_gopher_original_1024x504.jpg   --output gopher_500x500.jpg
curl http://localhost:8080/fill/1024/252/imagesource/_gopher_original_1024x504.jpg  --output gopher_1024x252.jpg
curl http://localhost:8080/fill/2000/1000/imagesource/_gopher_original_1024x504.jpg --output gopher_2000x1000.jpg
