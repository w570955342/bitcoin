#!/bin/bash

./block send  1B7x744vTRR6EYQL3VAnfPYVaqdQaRZSy8 14oXWSDaHQ3vYyfJJRSk1DziRiFF7uqyRN 10 1Agbes4R88GnNpouXiim2PyouqyA3c2kp2 "张三转李四10"
./block send  1B7x744vTRR6EYQL3VAnfPYVaqdQaRZSy8 15cVx9uD5qhistnzNTn2tsyzkb7KWkrtLq 20 1Agbes4R88GnNpouXiim2PyouqyA3c2kp2 "张三转王五20"
./checkBalance.sh

echo "========================================================================================"
./block send  15cVx9uD5qhistnzNTn2tsyzkb7KWkrtLq 14oXWSDaHQ3vYyfJJRSk1DziRiFF7uqyRN 2 1Agbes4R88GnNpouXiim2PyouqyA3c2kp2 "王五转李四2"
./block send  15cVx9uD5qhistnzNTn2tsyzkb7KWkrtLq 14oXWSDaHQ3vYyfJJRSk1DziRiFF7uqyRN 3 1Agbes4R88GnNpouXiim2PyouqyA3c2kp2 "王五转李四3"
./block send  15cVx9uD5qhistnzNTn2tsyzkb7KWkrtLq 1B7x744vTRR6EYQL3VAnfPYVaqdQaRZSy8 5 1Agbes4R88GnNpouXiim2PyouqyA3c2kp2 "王五转张三5"
./checkBalance.sh

echo "========================================================================================"
./block send  14oXWSDaHQ3vYyfJJRSk1DziRiFF7uqyRN 1Vw2BqqhhZmZGnCdcrMEVBPuXJcaSs2VG 14 1Agbes4R88GnNpouXiim2PyouqyA3c2kp2 "李四转赵六14"
./checkBalance.sh

