#!/bin/bash

sudo journalctl -u chainservice --since "1 minute ago" -f
