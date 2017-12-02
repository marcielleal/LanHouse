#!/bin/bash

((i=0))
while ((i<$1)); do
	echo $i
	./LanHouse > outputs/output$i
	((i+=1))
	echo $i
done
