#!/bin/bash

docker build -t armory-backend ./backend 

docker build -t armory-nginx ./nginx 