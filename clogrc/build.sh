#!/usr/bin/env bash
# clog>  build # build & inject metadata into clog
# extra> push main executables into tmp/
#                             _                            _     _                         _          _       __         _      _
#   ___   _ __   ___   _ _   | |_   ___  __ _   ___   __  | |_  | |  ___  __ __ __  __ _  | |_   __  | |_    / _|  ___  | |  __| |  ___   _ _
#  / _ \ | '_ \ / -_) | ' \  |  _| (_-< / _` | |___| / _| |  _| | | |___| \ V  V / / _` | |  _| / _| | ' \  |  _| / _ \ | | / _` | / -_) | '_|
#  \___/ | .__/ \___| |_||_|  \__| /__/ \__, |       \__|  \__| |_|        \_/\_/  \__,_|  \__| \__| |_||_| |_|   \___/ |_| \__,_| \___| |_|
#        |_|                            |___/
# ------------------------------------------------------------------------------
# load build config and script helpers
eval "$(clog Source project config)"                            # configure project
eval "$(clog Inc)"                                       # include clog helpers (sh, zsh & bash)
eval "$(clog Cat core/sh/help-golang.sh)"

# highlight colors
cLnx="$cC";cMac="$cW";cWin="$cE";cArm="$cS";cAmd="$cH"

fInfo "Building Project$cS $PROJECT $cT using$cC clog Cat$cF core/sh/help-golang.sh"

#clog Check
[ $? -gt 0 ] && exit 1
# ------------------------------------------------------------------------------

# ensure tmp dir exists
mkdir -p tmp

branch="$(clog git branch)"
hash="$(clog git hash head)"                                            # use the head hash as the build hash
suffix="" && [[ "$branch" != "main" ]] && suffix="$branch"              # use the branch name as the suffix
app=opentsg-ctl-watchfolder                                             # command you type to run the build
title="OpenTSG Render Node"                                             # title of the software
linkerPath="github.com/opentsg/opentsg-ctl-watchfolder/semver.SemVerInfo" # go tool objdump -S tmp/opentsg-ctl-watchfolder-amd-lnx|grep /semver.SemVerInfo

fGoBuild tmp/$app-amd-lnx     linux   amd64 $hash "$suffix" $app "$title" "$linkerPath"
fGoBuild tmp/$app-amd-win.exe windows amd64 $hash "$suffix" $app "$title" "$linkerPath"
fGoBuild tmp/$app-amd-mac     darwin  amd64 $hash "$suffix" $app "$title" "$linkerPath"
fGoBuild tmp/$app-arm-lnx     linux   arm64 $hash "$suffix" $app "$title" "$linkerPath"
fGoBuild tmp/$app-arm-win.exe windows arm64 $hash "$suffix" $app "$title" "$linkerPath"
fGoBuild tmp/$app-arm-mac     darwin  arm64 $hash "$suffix" $app "$title" "$linkerPath"

fInfo "${cT}All built to the$cF tmp/$cT folder\n"