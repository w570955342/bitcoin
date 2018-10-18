#!/bin/bash
rm block
rm *.db
go build -o block
./block