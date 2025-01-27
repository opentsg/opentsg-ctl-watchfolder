#                             _                       _          _
#   ___   _ __   ___   _ _   | |_   ___  __ _   ___  | |  __ _  | |__
#  / _ \ | '_ \ / -_) | ' \  |  _| (_-< / _` | |___| | | / _` | | '_ \
#  \___/ | .__/ \___| |_||_|  \__| /__/ \__, |       |_| \__,_| |_.__/
#        |_|                            |___/
# Create params to control the build
eval "$(clog Inc)"

export PROJECT="$(basename $(pwd))"
export vCODE=$(clog git vcode)
export vCodeType="Golang"
export bBASE="opentsg"
export bHASH="$(clog git hash head)"
export bMSG=$(clog git message latest)
# add a suffix to any build not on the main branch
export bSUFFIX="$(clog git branch)" && [[ "$bSUFFIX"=="main" ]] && bSUFFIX=""
# keep the private msg stuff in the mrmxf namespace. Only use opentsg for
# public stuff from now.
export bDOCKER_NS="mrmxf"
