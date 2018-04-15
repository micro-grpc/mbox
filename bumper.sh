#!/bin/bash

GO_OS=${GO_OS:-"linux"}

function detect_os {
    # Detect the OS name
    case "$(uname -s)" in
      Darwin)
        host_os=darwin
        ;;
      Linux)
        host_os=linux
        ;;
      *)
        echo "Unsupported host OS. Must be Linux or Mac OS X." >&2
        exit 1
        ;;
    esac

   GO_OS="${host_os}"
}

detect_os

#old_version=$(grep RELEASE?= ./Makefile | sed -e 's/RELEASE?=\(.*\)/\1/g')
old_version=$(cat ./RELEASE)
supposed_version=$(echo $old_version | awk -F. '{$NF = $NF + 1;} 1' | sed 's/ /./g')

echo "Current version $old_version."
echo -n "Please enter bumped version [$supposed_version]:"
read new_version
new_version=${new_version:-$supposed_version}

if [[ $new_version =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "Bumped version: $new_version"
else
    echo "Version is incorrect, please use X.X.X format (ie: 0.17.3)"
    exit
fi

if [ "${GO_OS}" == "darwin" ]; then
    echo $new_version > ./RELEASE
else
    echo $new_version > ./RELEASE
fi
