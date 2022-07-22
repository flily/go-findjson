#!/bin/bash

GO_VERSION_LIST=(
    "1.18"
    "1.17"
    "1.16"
    "1.15"
    "1.14"
    "1.13"
    "1.12"
)

GO_IMAGE_OS_LIST=(
    "alpine"
)

main()
{
    local wd="/go/src/go-findjson"
    for gover in "${GO_VERSION_LIST[@]}"; do
        for goos in "${GO_IMAGE_OS_LIST[@]}"; do
            local image_name="golang:${gover}-${goos}"
            echo "Checking on go ${gover} (${goos}): ${image_name}"

            docker run -v "$(pwd):${wd}" -w "${wd}" "${image_name}" go test -cover ./...
        done
    done
}

main "$@"
