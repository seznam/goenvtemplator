#!/bin/bash

RELEASE=${RELEASE:-jessie}
PRIORITY=${PRIORITY:-505}

echo -e "Package: *\nPin: release n=${RELEASE}-backports\nPin-Priority: ${PRIORITY}\n" >> /etc/apt/preferences.d/backports.pref
