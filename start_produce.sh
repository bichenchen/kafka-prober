#!/bin/bash
topic="kafka-prober-topic1"
mode="produce"
username="normal"
passwd="normal-secret-test"
cluster=$1
log="/home/work/kafka-prober-${cluster}.log"
cluster=`grep ${clustname} kafka-cluster.txt|cut -d'/' -f1`
if [ "${cluster}X" == "X" ];then
  echo "cluster name ${cluster} not exits"
  exit
fi
nohup ./kafka-prober -brokers "${brokers}" -topic "${topic}" -mode ${mode} -username="${username}" -passwd ${passwd} -cluster ${cluster} -n ${log} >kafka-prober-${mode}-${cluster}.exec.log &